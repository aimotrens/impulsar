package engine

import (
	"fmt"
	"os"

	"github.com/aimotrens/impulsar/internal/model"
)

func (e *Engine) RunJob(job string) {
	if scheduledJob, ok := e.GetJob(job); ok {
		e.executeJob(scheduledJob)
		return
	}

	fmt.Printf("Job %s not found\n", job)
}

func (e *Engine) executeJob(j *model.Job) {
	e.evaluateConditionalField(j)

	exitWithFinalizer := func() {
		_ = e.runScriptBlock(j, j.ScriptFinalize, "via script:finalize")

		for _, finalize := range j.JobsFinalize {
			e.RunJob(finalize)
		}

		os.Exit(1)
	}

	if e.evaluateIfCondition(j) {
		for _, pre := range j.JobsPre {
			e.RunJob(pre)
		}

		if j.ScriptPre != nil {
			if err := e.runScriptBlock(j, j.ScriptPre, "via script:pre"); err != nil {
				exitWithFinalizer()
			}
		}

		if j.Foreach != nil {
			for _, f := range j.Foreach {
				for k, v := range f {
					e.Variables[k] = v
				}

				if err := e.runScriptBlock(j, j.Script, "via foreach"); err != nil {
					exitWithFinalizer()
				}
			}
		} else {
			if err := e.runScriptBlock(j, j.Script, ""); err != nil {
				exitWithFinalizer()
			}
		}

		if j.ScriptPost != nil {
			if err := e.runScriptBlock(j, j.ScriptPost, "via script:post"); err != nil {
				exitWithFinalizer()
			}
		}

		if j.ScriptFinalize != nil {
			_ = e.runScriptBlock(j, j.ScriptFinalize, "view script:finalize")
		}

		for _, post := range j.JobsPost {
			e.RunJob(post)
		}

		for _, fin := range j.JobsFinalize {
			e.RunJob(fin)
		}
	} else {
		fmt.Printf("Job %s skipped, no condition matched\n", j.Name)
	}
}

func (e *Engine) execCommand(j *model.Job, script string) error {
	if shell, ok := e.shellMap[j.Shell.Type]; !ok {
		fmt.Printf("Shell type %s not supported\n", j.Shell.Type)
		os.Exit(1)
		return nil
	} else {
		return shell.Execute(j, script)
	}
}

func (e *Engine) evaluateIfCondition(j *model.Job) bool {
	if j.If == nil {
		return true
	}

	envVars := e.aggregateEnvVars(j)

	for _, varSet := range j.If {
		var success = true

		for k, v := range varSet {
			success = success && (envVars[k] == v)
		}

		if success {
			return true
		}
	}

	return false
}

func (e *Engine) evaluateConditionalField(j *model.Job) {
	if j.Conditional == nil {
		return
	}

	var envVars = e.aggregateEnvVars(j)

	for _, v := range j.Conditional {
		for _, varSet := range v.If {
			var success = true

			for k, v := range varSet {
				success = success && (envVars[k] == v)
			}

			if success {
				if v.Overwrite != nil {
					j.Overwrite(v.Overwrite)
				}

				return
			}
		}
	}
}
