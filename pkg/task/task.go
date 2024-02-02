package task

import (
	"fmt"
	"time"

	"sandboxer/pkg/logging"
	"sandboxer/pkg/state"
)

type ID int64

/*
var (

	count ID

)

	func TaskNumber() ID {
		return (ID)(atomic.AddInt64((*int64)(&count), 1))
	}
*/
type Task struct {
	//	mx         sync.Mutex
	Number     ID
	SubmitTime time.Time
	Path       string
	State      state.State
	Message    string
	SandboxID  string
}

func NewTask(id ID, path string) *Task {
	return &Task{
		Number:     id,
		SubmitTime: time.Now(),
		Path:       path,
		State:      state.StateNew,
		Message:    "",
		SandboxID:  "",
	}
}

//func (t *Task) lockUnlock() func() {
//		t.mx.Lock()
//		return t.mx.Unlock
//}

func (t *Task) SetState(newState state.State) {
	logging.Debugf("SetState(%v)", newState)
	t.State = newState
	t.Message = ""
}

func (t *Task) SetID(id string) {
	t.SandboxID = id
}

func (t *Task) VOneID() string {
	return t.SandboxID
}

func (t *Task) SetSandboxID(sandboxID string) {
	t.SandboxID = sandboxID
}

func (t *Task) String() string {
	return fmt.Sprintf("Task %d; submitted on: %v; state: %v; id: %s; message: %s, path: %s", t.Number, t.SubmitTime, t.State, t.SandboxID, t.Message, t.Path)
}

func (t *Task) SetError(err error) {
	t.State = state.StateError
	t.Message = err.Error()
}

func (t *Task) SetMessage(message string) {
	t.Message = message
}
