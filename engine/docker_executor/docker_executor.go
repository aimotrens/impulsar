package dockerexecutor

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/aimotrens/impulsar/engine"
	"github.com/aimotrens/impulsar/model"
)

type DockerExecutor struct {
	*engine.Engine
}

func init() {
	engine.RegisterExecutor(model.SHELL_TYPE_DOCKER, func(e *engine.Engine) engine.Shell {
		return &DockerExecutor{Engine: e}
	})
}

func (e *DockerExecutor) Execute(j *model.Job, script string) {
	wd := j.WorkDir
	if !filepath.IsAbs(wd) {
		currentWorkDir, _ := os.Getwd()
		wd = path.Join(currentWorkDir, wd)
	}

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
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	for key, value := range j.Variables {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	scriptExpanded := os.Expand(script, e.LookupVarFunc(j))

	args = append(args, "--entrypoint", j.Shell.BootCommand[0], j.Shell.Image)
	args = append(args, j.Shell.BootCommand[1:]...)
	args = append(args, scriptExpanded)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = &engine.JobOutputUnifier{Job: j, ScriptLine: &script, Writer: os.Stdout}
	cmd.Stderr = &engine.JobOutputUnifier{Job: j, ScriptLine: &script, Writer: os.Stderr}
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Command %s failed\n%s\n", script, err)

		if !j.AllowFail {
			os.Exit(1)
		}
	}
}
