package task

import (
	"examen/pkg/state"
	"time"
)

type Task struct {
	submitTime time.Time
	path       string
	state      state.State
	id         string
}

func NewTask(path string) *Task {
	return &Task{
		submitTime: time.Now(),
		path:       path,
		state:      state.StateUnknown,
		id:         "",
	}
}

func (t *Task) Upload() {

}
