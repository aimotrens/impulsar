package engine

import (
	"io"
	"os"
	"strings"

	"github.com/aimotrens/impulsar/internal/model"
)

type jobOutputCorrector struct {
	job          *model.Job
	outputTarget io.Writer
}

func (s *jobOutputCorrector) Write(p []byte) (n int, err error) {
	n = len(p)

	tmpOutput := string(p)

	if s.job.Shell.Type == model.SHELL_TYPE_BASH {
		tmpOutput = strings.ReplaceAll(tmpOutput, "\n", "\r\n")
	}

	_, err = s.outputTarget.Write([]byte(tmpOutput))
	return
}

func GetCmdOutputTarget(j *model.Job) (out, err io.Writer) {
	return &jobOutputCorrector{job: j, outputTarget: os.Stdout}, &jobOutputCorrector{job: j, outputTarget: os.Stderr}
}
