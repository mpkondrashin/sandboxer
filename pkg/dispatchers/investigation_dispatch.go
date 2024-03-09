/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

investigation_dispatch.go

Download investigation package
*/
package dispatchers

import (
	"sandboxer/pkg/task"
)

type InvestigationDispatch struct {
	BaseDispatcher
}

func NewInvestigationDispatch(d BaseDispatcher) *InvestigationDispatch {
	return &InvestigationDispatch{
		BaseDispatcher: d,
	}
}

func (d *InvestigationDispatch) InboundChannel() task.Channel {
	return task.ChInvestigation
}

func (d *InvestigationDispatch) ProcessTask(tsk *task.Task) error {
	sbox, err := d.Sandbox()
	if err != nil {
		return err
	}
	zipFilePath, err := tsk.InvestigationPath()
	if err != nil {
		return err
	}
	if err := sbox.GetInvestigation(tsk.SandboxID, zipFilePath); err != nil {
		return err
	}
	tsk.SetInvestigation(zipFilePath)
	tsk.SetChannel(task.ChDone)
	d.list.Updated()
	return nil
}
