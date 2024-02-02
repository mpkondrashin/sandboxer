package dispatchers

import (
	"context"
	"fmt"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/state"
	"sandboxer/pkg/task"
	"time"

	"github.com/mpkondrashin/vone"
)

type WaitDispatch struct {
	BaseDispatcher
}

func NewWaitDispatch(d BaseDispatcher) *WaitDispatch {
	return &WaitDispatch{
		BaseDispatcher: d,
	}
}

func (*WaitDispatch) InboundChannel() int {
	return ChWait
}

func (d *WaitDispatch) ProcessTask(tsk *task.Task) error {
	tsk.SetState(state.StateCheck)
	vOne, err := d.vOne()
	if err != nil {
		return err
	}
	logging.Debugf("id: %s", tsk.SandboxID)
	status, err := vOne.SandboxSubmissionStatus(tsk.SandboxID).Do(context.TODO())
	if err != nil {
		return fmt.Errorf("SandboxSubmissionStatus: %w", err)
	}
	logging.Debugf("%s Status: %v", tsk.SandboxID, status.Status)
	switch status.Status {
	case vone.StatusSucceeded:
		tsk.SetState(state.StateWaitForResult)
		d.list.Updated()
		d.Channel(ChResult) <- tsk.Number
		return nil
	case vone.StatusRunning:
		tsk.SetState(state.StateInspected)
		d.list.Updated()
		time.Sleep(d.conf.VisionOne.Sleep)
		d.Channel(ChWait) <- tsk.Number
	case vone.StatusFailed:
		return fmt.Errorf("%s: %s", status.Error.Code, status.Error.Message)
	default:
		return fmt.Errorf("unknown status: %s", status.Status)
	}
	return nil
}
