package engine

import (
	"bufio"
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

func (e *Engine) CollectArgs(job string) {
	if j, ok := e.impulsar[job]; ok {
		readArgs(e, j)

		for _, pre := range j.JobsPre {
			e.CollectArgs(pre)
		}

		for _, post := range j.JobsPost {
			e.CollectArgs(post)
		}

		return
	}

	fmt.Printf("Job %s not found\n", job)
}

func (e *Engine) RunJob(job string) {
	if scheduledJob, ok := e.impulsar[job]; ok {
		e.executeJob(scheduledJob)
		return
	}

	fmt.Printf("Job %s not found\n", job)
}

func (e *Engine) executeJob(j *model.Job) {
	readArgs(e, j)
	evaluateConditionalField(e, j)

	runScriptBlock := func(isForeach bool) {
		for _, script := range j.Script {
			if script == "STOP" {
				fmt.Printf("Job %s failed, due to STOP command\n", j.Name)
				os.Exit(1)
			}

			foreachSuffix := ""
			if isForeach {
				foreachSuffix = " via foreach"
			}

			fmt.Printf("[%s] (%s)%s\n",
				j.Name,
				strings.ReplaceAll(
					strings.Trim(script, "\n"),
					"\n",
					"; "),
				foreachSuffix,
			)

			e.execCommand(j, script)
		}
	}

	if evaluateIfCondition(e, j) {
		for _, pre := range j.JobsPre {
			e.RunJob(pre)
		}

		if j.Foreach != nil {
			for _, f := range j.Foreach {
				for k, v := range f {
					e.variables[k] = v
				}

				runScriptBlock(true)
			}
		} else {
			runScriptBlock(false)
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
	case model.SHELL_TYPE_SSH:
		e.execSshCommand(j, script)
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

	envVars["os"] = runtime.GOOS
	envVars["arch"] = runtime.GOARCH

	return envVars
}

func readArgs(e *Engine, j *model.Job) {
	for arg, description := range j.Arguments {
		if _, ok := j.Variables[arg]; ok {
			continue
		}

		if _, ok := e.variables[arg]; ok {
			continue
		}

		if val, ok := os.LookupEnv(arg); ok {
			j.Variables[arg] = val
			continue
		}

		fmt.Printf("[%s] %s (%s): ", j.Name, arg, description)

		var value string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			value = scanner.Text()
		}

		j.Variables[arg] = value
	}
}
