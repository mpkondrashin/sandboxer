/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

wait_dispatch.go

Wait for inspection result
*/
package dispatchers

import (
	"context"
	"fmt"
	"sandboxer/pkg/logging"
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

func (*WaitDispatch) InboundChannel() task.Channel {
	return task.ChWait
}

func (d *WaitDispatch) ProcessTask(tsk *task.Task) error {
	//tsk.SetState(task.StateCheck)

	vOne, err := d.vOne()
	if err != nil {
		return err
	}
	logging.Debugf("id: %s", tsk.SandboxID)
	status, err := vOne.SandboxSubmissionStatus(tsk.SandboxID).Do(context.TODO())
	if err != nil {
		return fmt.Errorf("SandboxSubmissionStatus: %w", err)
	}
	tsk.Deactivate()
	logging.Debugf("%s Status: %v", tsk.SandboxID, status.Status)
	switch status.Status {
	case vone.StatusSucceeded:
		//tsk.SetState(task.StateWaitForResult)
		//d.list.Updated()
		tsk.SetChannel(task.ChResult)
		return nil
	case vone.StatusRunning:
		//		tsk.SetState(task.StateInspected)
		//		d.list.Updated()
		time.Sleep(d.conf.VisionOne.Sleep)
		tsk.SetChannel(task.ChWait)
	case vone.StatusFailed:
		if status.Error.Code == "Unsupported" {
			tsk.SetChannel(task.ChDone)
			tsk.SetRiskLevel(task.RiskLevelUnsupported)
			tsk.SetMessage(status.Error.Message)
			d.list.Updated()
			return nil
		}
		return fmt.Errorf("%s: %s", status.Error.Code, status.Error.Message)
	default:
		return fmt.Errorf("unknown status: %s", status.Status)
	}
	return nil
}
