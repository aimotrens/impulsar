package engine

import (
	"fmt"
	"io"
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
