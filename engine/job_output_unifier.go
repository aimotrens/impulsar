package engine

import (
	"io"
	"runtime"
	"strings"

	"github.com/aimotrens/impulsar/model"
)

type JobOutputUnifier struct {
	Job        *model.Job
	ScriptLine *string
	Writer     io.Writer
}

func (s *JobOutputUnifier) Write(p []byte) (n int, err error) {
	n = len(p)

	tmpOutput := string(p)

	if runtime.GOOS == "windows" && s.Job.Shell.Type == model.SHELL_TYPE_BASH {
		tmpOutput = strings.ReplaceAll(tmpOutput, "\n", "\r\n")
	}

	_, err = s.Writer.Write([]byte(tmpOutput))
	return
}
