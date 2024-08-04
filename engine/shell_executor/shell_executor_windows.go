package shellexecutor

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/aimotrens/impulsar/model"
)

func (e *ShellExecutor) prepareCmdEnv(j *model.Job, cmd *exec.Cmd) {
	for key, value := range e.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	for key, value := range j.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	if j.Shell.Type == model.SHELL_TYPE_BASH {
		var wslEnv []string

		for key := range e.Variables {
			if key == "PATH" {
				continue
			}

			if j.VariablesExcluded != nil && slices.Contains(j.VariablesExcluded, key) {
				continue
			}

			wslEnv = append(wslEnv, key+"/u")
		}

		for key := range j.Variables {
			wslEnv = append(wslEnv, key+"/u")
		}

		cmd.Env = append(cmd.Env, "WSLENV="+strings.Join(wslEnv, ":"))
	}
}
