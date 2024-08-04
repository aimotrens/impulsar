//go:build !windows

package shellexecutor

import (
	"fmt"
	"os/exec"

	"github.com/aimotrens/impulsar/model"
)

func (e *ShellExecutor) prepareCmdEnv(j *model.Job, cmd *exec.Cmd) {
	for key, value := range e.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	for key, value := range j.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
}
