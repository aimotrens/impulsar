package engine

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/aimotrens/impulsar/model"
	"github.com/dop251/goja"
)

type Engine struct {
	impulsar  model.ImpulsarList
	variables map[string]string
}

func New(d model.ImpulsarList, variables map[string]string) *Engine {
	return &Engine{
		impulsar:  d,
		variables: variables,
	}
}

func (e *Engine) RunJob(job string) {
	if scheduledJob, ok := e.impulsar[job]; !ok {
		fmt.Printf("Job %s not found\n", job)
		return
	} else {
		e.executeJob(scheduledJob)
	}
}

func (e *Engine) executeJob(j *model.Job) {
	evaluateConditionalField(e, j)

	if evaluateIfCondition(e, j) {
		for _, pre := range j.JobsPre {
			e.RunJob(pre)
		}

		for _, script := range j.Script {
			if script == "STOP" {
				fmt.Printf("Job %s failed, due to STOP command\n", j.Name)
				os.Exit(1)
			}

			fmt.Printf("[%s] (%s)\n",
				j.Name,
				strings.ReplaceAll(
					strings.Trim(script, "\n"),
					"\n",
					"; "),
			)

			e.execCommand(j, script)
		}

		for _, post := range j.JobsPost {
			e.RunJob(post)
		}
	} else {
		fmt.Printf("Job %s skipped, no condition matched\n", j.Name)
	}
}

func (e *Engine) execCommand(j *model.Job, script string) {
	switch j.Shell.Type {
	case model.SHELL_TYPE_POWERSHELL, model.SHELL_TYPE_PWSH, model.SHELL_TYPE_BASH, model.SHELL_TYPE_CUSTOM:
		e.execShellCommand(j, script)
	case model.SHELL_TYPE_DOCKER:
		e.execDockerCommand(j, script)
	}
}

func variableMapper(e *Engine, j *model.Job) func(string) string {
	return func(s string) string {
		if v, ok := e.variables[s]; ok {
			return v
		}

		if v, ok := j.Variables[s]; ok {
			return v
		}

		return ""
	}
}

func evaluateIfCondition(e *Engine, j *model.Job) bool {
	if j.If == nil {
		return true
	}

	var envVars = collectEnvVars(e, j)

	vm := goja.New()

	vm.Set("env", envVars)

	for _, v := range j.If {
		res, _ := vm.RunString(v)

		if res.ToBoolean() {
			return true
		}
	}

	return false
}

func evaluateConditionalField(e *Engine, j *model.Job) {
	if j.Conditional == nil {
		return
	}

	var envVars = collectEnvVars(e, j)

	vm := goja.New()
	vm.Set("env", envVars)

	for _, v := range j.Conditional {
		for _, ifCondition := range v.If {
			res, _ := vm.RunString(ifCondition)

			if res.ToBoolean() {
				if v.Overwrite != nil {
					j.Overwrite(v.Overwrite)
				}

				return
			}
		}
	}
}

func collectEnvVars(e *Engine, j *model.Job) map[string]string {
	var envVars = make(map[string]string)

	envVars["os"] = runtime.GOOS
	envVars["arch"] = runtime.GOARCH

	for _, v := range os.Environ() {
		var kv = strings.Split(v, "=")
		envVars[strings.ToLower(kv[0])] = kv[1]
	}

	for key, value := range e.variables {
		envVars[strings.ToLower(key)] = value
	}

	for key, value := range j.Variables {
		envVars[strings.ToLower(key)] = value
	}

	return envVars
}
