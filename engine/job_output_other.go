//go:build !windows

package engine

import (
	"io"
	"os"

	"github.com/aimotrens/impulsar/model"
)

func GetCmdOutputTarget(j *model.Job) (out, err io.Writer) {
	return os.Stdout, os.Stderr
}
