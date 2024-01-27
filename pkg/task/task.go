package task

import (
	"fmt"
	"sync/atomic"
	"time"

	"sandboxer/pkg/state"
)

type ID int64

var (
	count ID
)

func TaskNumber() ID {
	return (ID)(atomic.AddInt64((*int64)(&count), 1))
}

type Task struct {
	Number     ID
	SubmitTime time.Time
	Path       string
	State      state.State
	Message    string
	SandboxID  string
}

func NewTask(path string) *Task {
	return &Task{
		Number:     TaskNumber(),
		SubmitTime: time.Now(),
		Path:       path,
		State:      state.StateNew,
		Message:    "",
		SandboxID:  "",
	}
}

func (t *Task) SetState(newState state.State) {
	t.State = newState
}
func (t *Task) SetID(id string) {
	t.SandboxID = id
}

func (t *Task) VOneID() string {
	return t.SandboxID
}

func (t *Task) String() string {
	return fmt.Sprintf("Task %d; submitted on: %v; state: %v; id: %s; path: %s", t.Number, t.SubmitTime, t.State, t.SandboxID, t.Path)
}

func (t *Task) SetError(err error) {
	t.State = state.StateError
	t.Message = err.Error()
}

func (t *Task) SetMessage(message string) {
	t.Message = message
}
