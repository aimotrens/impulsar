package engine

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aimotrens/impulsar/model"
)

type jobOutputPrefixer struct {
	Job           *model.Job
	ScriptLine    *string
	Writer        io.Writer
	prefixWritten bool
}

func (s *jobOutputPrefixer) Write(p []byte) (n int, err error) {
	n = len(p)

	tmpOutput := string(p)
	if !s.prefixWritten {
		tmpOutput = fmt.Sprintf("[%s] (%s)\n%s",
			s.Job.Name,
			strings.ReplaceAll(
				strings.Trim(*s.ScriptLine, "\n"),
				"\n",
				"; "),
			p,
		)
		s.prefixWritten = true
	}

	if runtime.GOOS == "windows" && s.Job.Shell.Type == model.SHELL_TYPE_BASH {
		tmpOutput = strings.ReplaceAll(tmpOutput, "\n", "\r\n")
	}

	_, err = s.Writer.Write([]byte(tmpOutput))
	return
}

type Engine struct {
	impulsar  model.ImpulsarList
	variables map[string]string
}

func New(d model.ImpulsarList, variables map[string]string) *Engine {
	return &Engine{
		impulsar:  d,
		variables: variables,
	}
}

func (e *Engine) RunJob(job string) {
	if scheduledJob, ok := e.impulsar[job]; !ok {
		fmt.Printf("Job %s not found\n", job)
		return
	} else {
		e.executeJob(scheduledJob)
	}
}

func (e *Engine) executeJob(j *model.Job) {
	for _, pre := range j.JobsPre {
		e.RunJob(pre)
	}

	for _, script := range j.Script {
		e.execCommand(j, script)
	}

	for _, post := range j.JobsPost {
		e.RunJob(post)
	}
}

func (e *Engine) execCommand(j *model.Job, script string) {
	switch j.Shell.Type {
	case model.SHELL_TYPE_POWERSHELL, model.SHELL_TYPE_PWSH, model.SHELL_TYPE_BASH, model.SHELL_TYPE_CUSTOM:
		e.execShellCommand(j, script)
	case model.SHELL_TYPE_DOCKER:
		e.execDockerCommand(j, script)
	}
}

func (e *Engine) execShellCommand(j *model.Job, script string) {
	cmd := exec.Command(j.Shell.BootCommand[0], append(j.Shell.BootCommand[1:], script)...)
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

func (e *Engine) execDockerCommand(j *model.Job, script string) {
	wd := j.WorkDir
	if !filepath.IsAbs(wd) {
		currentWorkDir, _ := os.Getwd()
		wd = path.Join(currentWorkDir, wd)
	}

	currentWorkDir, _ := os.Getwd()
	args := []string{"run", "--rm", "-v", currentWorkDir + ":/workdir", "--workdir=" + path.Join("/workdir", j.WorkDir)}

	if runtime.GOOS != "windows" {
		args = append(args, "--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
	}

	for _, v := range e.variables {
		args = append(args, "-e", v)
	}

	for key, value := range j.Variables {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, "--entrypoint", j.Shell.BootCommand[0], j.Shell.Image)
	args = append(args, j.Shell.BootCommand[1:]...)
	args = append(args, script)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = &jobOutputPrefixer{Job: j, ScriptLine: &script, Writer: os.Stdout}
	cmd.Stderr = &jobOutputPrefixer{Job: j, ScriptLine: &script, Writer: os.Stderr}
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Command %s failed\n%s\n", script, err)

		if !j.AllowFail {
			os.Exit(1)
		}
	}
}
