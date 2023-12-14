package engine

import (
	"fmt"

	"github.com/aimotrens/impulsar/model"
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
	for _, pre := range j.JobsPre {
		e.RunJob(pre)
	}

	for _, script := range j.Script {
		e.execCommand(j, script)
	}

	for _, post := range j.JobsPost {
		e.RunJob(post)
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
