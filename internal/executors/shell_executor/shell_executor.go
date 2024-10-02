package shellexecutor

import (
	"fmt"
	"os/exec"

	"github.com/aimotrens/impulsar/internal/engine"
	"github.com/aimotrens/impulsar/internal/model"
)

type ShellExecutor struct {
	*engine.Engine
}

func init() {
	constructor := func(e *engine.Engine) engine.Executor {
		return &ShellExecutor{Engine: e}
	}

	engine.RegisterExecutor(model.SHELL_TYPE_BASH, constructor)
	engine.RegisterExecutor(model.SHELL_TYPE_POWERSHELL, constructor)
	engine.RegisterExecutor(model.SHELL_TYPE_PWSH, constructor)
	engine.RegisterExecutor(model.SHELL_TYPE_CUSTOM, constructor)
}

func (e *ShellExecutor) Execute(j *model.Job, script string) error {
	scriptExpanded := e.ExpandVarsWithTemplateEngine(script, j)

	cmd := exec.Command(j.Shell.BootCommand[0], append(j.Shell.BootCommand[1:], scriptExpanded)...)
	cmd.Stdout, cmd.Stderr = engine.GetCmdOutputTarget(j)

	e.prepareCmdEnv(j, cmd)

	cmd.Dir = j.WorkDir

	err := runCmd(cmd, j.Shell.Interactive)

	if err != nil {
		fmt.Printf("command %s failed:\n%s\n", script, err)

		if !j.AllowFail {
			return fmt.Errorf("command %s failed:\n%s", script, err)
		}
	}

	return nil
}
