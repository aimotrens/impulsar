package shellexecutor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/aimotrens/impulsar/engine"
	"github.com/aimotrens/impulsar/model"
)

type ShellExecutor struct {
	*engine.Engine
}

func init() {
	constructor := func(e *engine.Engine) engine.Shell {
		return &ShellExecutor{Engine: e}
	}

	engine.RegisterExecutor(model.SHELL_TYPE_BASH, constructor)
	engine.RegisterExecutor(model.SHELL_TYPE_POWERSHELL, constructor)
	engine.RegisterExecutor(model.SHELL_TYPE_PWSH, constructor)
	engine.RegisterExecutor(model.SHELL_TYPE_CUSTOM, constructor)
}

func (e *ShellExecutor) Execute(j *model.Job, script string) error {
	scriptExpanded := e.ExpandVarsWithTemplateEngine(script, j)
	scriptExpanded = os.Expand(scriptExpanded, e.LookupVarFunc(j))

	cmd := exec.Command(j.Shell.BootCommand[0], append(j.Shell.BootCommand[1:], scriptExpanded)...)
	cmd.Stdout = &engine.JobOutputUnifier{Job: j, ScriptLine: &script, Writer: os.Stdout}
	cmd.Stderr = &engine.JobOutputUnifier{Job: j, ScriptLine: &script, Writer: os.Stderr}

	for key, value := range e.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	for key, value := range j.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	if runtime.GOOS == "windows" && j.Shell.Type == model.SHELL_TYPE_BASH {
		var wslEnv []string

		for key := range e.Variables {
			if key == "PATH" {
				continue
			}

			wslEnv = append(wslEnv, key+"/u")
		}

		for key := range j.Variables {
			wslEnv = append(wslEnv, key+"/u")
		}

		cmd.Env = append(cmd.Env, "WSLENV="+strings.Join(wslEnv, ":"))
	}

	cmd.Dir = j.WorkDir

	err := cmd.Run()

	if err != nil {
		fmt.Printf("Command %s failed:\n%s\n", script, err)

		if !j.AllowFail {
			return fmt.Errorf("Command %s failed:\n%s\n", script, err)
		}
	}

	return nil
}
