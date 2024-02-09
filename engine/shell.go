package engine

import "github.com/aimotrens/impulsar/model"

type Shell interface {
	Execute(j *model.Job, script string) error
}
