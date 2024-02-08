package task

import (
	"fmt"
	"time"

	"sandboxer/pkg/logging"
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
	State      State
	RiskLevel  RiskLevel
	Message    string
	SandboxID  string
	MD5        string
	SHA1       string
	SHA256     string
	Report     string
}

func NewTask(id ID, path string) *Task {
	return &Task{
		Number:     id,
		SubmitTime: time.Now(),
		Path:       path,
		State:      StateNew,
		RiskLevel:  RiskLevelUnknown,
		Message:    "",
		SandboxID:  "",
	}
}

//func (t *Task) lockUnlock() func() {
//		t.mx.Lock()
//		return t.mx.Unlock
//}

func (t *Task) SetState(newState State) {
	logging.Debugf("SetState(%v)", newState)
	t.State = newState
}

func (t *Task) GetState() string {
	if t.State == StateDone {
		return t.RiskLevel.String()
	}
	return t.State.String()
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
func (t *Task) SetRiskLevel(riskLevel RiskLevel) {
	t.State = StateDone
	t.RiskLevel = riskLevel
}

func (t *Task) SetError(err error) {
	t.State = StateDone
	t.RiskLevel = RiskLevelError
	t.Message = err.Error()
}

func (t *Task) SetMessage(message string) {
	t.Message = message
}

func (t *Task) SetReport(report string) {
	t.Report = report
}

func (t *Task) SetDigest(MD5, SHA1, SHA256 string) {
	if MD5 != "" {
		t.MD5 = MD5
	}
	if SHA1 != "" {
		t.SHA1 = SHA1
	}
	if SHA256 != "" {
		t.SHA256 = SHA256
	}
}
