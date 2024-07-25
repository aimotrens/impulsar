package engine

import (
	"runtime"

	"github.com/aimotrens/impulsar/model"
)

func (e *Engine) LookupVarFunc(j *model.Job) func(string) string {
	return func(s string) string {
		if v, ok := j.Variables[s]; ok {
			return v
		}

		if v, ok := e.Variables[s]; ok {
			return v
		}

		return ""
	}
}

func (e *Engine) aggregateEnvVars(j *model.Job) model.VariableMap {
	var envVars = make(model.VariableMap)

	for key, value := range e.Variables {
		envVars[key] = value
	}

	for key, value := range j.Variables {
		envVars[key] = value
	}

	envVars["OS"] = runtime.GOOS
	envVars["ARCH"] = runtime.GOARCH

	return envVars
}
