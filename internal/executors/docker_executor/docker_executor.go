package dockerexecutor

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"slices"

	"github.com/aimotrens/impulsar/internal/engine"
	"github.com/aimotrens/impulsar/internal/model"
)

type DockerExecutor struct {
	*engine.Engine
}

func init() {
	engine.RegisterExecutor(model.SHELL_TYPE_DOCKER, func(e *engine.Engine) engine.Executor {
		return &DockerExecutor{Engine: e}
	})
}

func (e *DockerExecutor) Execute(j *model.Job, script string) error {
	currentWorkDir, _ := os.Getwd()
	args := []string{"run", "--rm", "-v", currentWorkDir + ":/workdir", "--workdir=" + path.Join("/workdir", j.WorkDir)}

	if j.Shell.UidGid != "" {
		args = append(args, "--user", j.Shell.UidGid)
	} else {
		if runtime.GOOS == "windows" {
			args = append(args, "--user", "1000:1000")
		} else {
			args = append(args, "--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
		}
	}

	for key, value := range e.Variables {
		if key == "PATH" {
			continue
		}

		if j.VariablesExcluded != nil && slices.Contains(j.VariablesExcluded, key) {
			continue
		}

		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	for key, value := range j.Variables {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	scriptExpanded := e.ExpandVarsWithTemplateEngine(script, j)

	args = append(args, "--entrypoint", j.Shell.BootCommand[0], j.Shell.Image)
	args = append(args, j.Shell.BootCommand[1:]...)
	args = append(args, scriptExpanded)

	cmd := exec.Command("docker", args...)
	cmd.Stdout, cmd.Stderr = engine.GetCmdOutputTarget(j)
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Command %s failed\n%s\n", script, err)

		if !j.AllowFail {
			return fmt.Errorf("command %s failed\n%s", script, err)
		}
	}

	return nil
}
