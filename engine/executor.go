package engine

import "github.com/aimotrens/impulsar/model"

type Executor interface {
	Execute(j *model.Job, script string) error
}
