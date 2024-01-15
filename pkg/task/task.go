package task

import (
	"examen/pkg/state"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	count int64
)

func TaskNumber() int64 {
	return atomic.AddInt64(&count, 1)
}

type Task struct {
	Number     int64
	SubmitTime time.Time
	Path       string
	State      state.State
	Message    string
	vOneID     string
}

func NewTask(path string) *Task {
	return &Task{
		Number:     TaskNumber(),
		SubmitTime: time.Now(),
		Path:       path,
		State:      state.StateNew,
		Message:    "",
		vOneID:     "",
	}
}

func (t *Task) SetState(newState state.State) {
	t.State = newState
}
func (t *Task) SetID(id string) {
	t.vOneID = id
}

func (t *Task) VOneID() string {
	return t.vOneID
}

func (t *Task) String() string {
	return fmt.Sprintf("Task %d; submitted on: %v; state: %v; id: %s; path: %s", t.Number, t.SubmitTime, t.State, t.vOneID, t.Path)
}

func (t *Task) SetError(errMessage string) {
	t.State = state.StateError
	t.Message = errMessage
}

func (t *Task) SetMessage(message string) {
	t.Message = message
}
