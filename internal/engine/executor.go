package engine

import "github.com/aimotrens/impulsar/internal/model"

type Executor interface {
	Execute(j *model.Job, script string) error
}
