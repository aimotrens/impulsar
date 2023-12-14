package engine

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/aimotrens/impulsar/model"
)

func (e *Engine) execShellCommand(j *model.Job, script string) {
	scriptExpanded := os.Expand(script, variableMapper(e, j))

	cmd := exec.Command(j.Shell.BootCommand[0], append(j.Shell.BootCommand[1:], scriptExpanded)...)
	cmd.Stdout = &jobOutputPrefixer{Job: j, ScriptLine: &script, Writer: os.Stdout}
	cmd.Stderr = &jobOutputPrefixer{Job: j, ScriptLine: &script, Writer: os.Stderr}
	cmd.Env = os.Environ()

	for key, value := range e.variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	for key, value := range j.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	if runtime.GOOS == "windows" && j.Shell.Type == model.SHELL_TYPE_BASH {
		var wslEnv []string

		for key := range e.variables {
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
			os.Exit(1)
		}
	}
}
