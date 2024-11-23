package engine

import (
	"os"
	"strings"

	"github.com/aimotrens/impulsar/model"
)

type (
	Engine struct {
		jobList   map[string]*model.Job
		shellMap  map[string]Executor
		Variables model.VariableMap
	}

	ExecutorConstructor func(*Engine) Executor
)

var executorMap = make(map[string]ExecutorConstructor)

func RegisterExecutor(name string, constructor ExecutorConstructor) {
	executorMap[name] = constructor
}

func New(jobList map[string]*model.Job, additionalEnvVars model.VariableMap) *Engine {
	envVars := make(model.VariableMap)

	// Aggregate environment and additional variables
	{
		for _, v := range os.Environ() {
			var kv = strings.SplitN(v, "=", 2)
			envVars[kv[0]] = kv[1]
		}

		for key, value := range additionalEnvVars {
			envVars[key] = value
		}
	}

	e := &Engine{
		jobList:   jobList,
		shellMap:  make(map[string]Executor),
		Variables: envVars,
	}

	for name, constructor := range executorMap {
		e.shellMap[name] = constructor(e)
	}

	return e
}

func (e *Engine) GetJob(name string) (j *model.Job, ok bool) {
	j, ok = e.jobList[name]
	return
}
