/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

upload_dispatch.go

Upload file for inspection
*/
package dispatchers

import (
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
)

type UploadDispatch struct {
	BaseDispatcher
}

func (*UploadDispatch) InboundChannel() task.Channel {
	return task.ChSubmit
}

func NewUploadDispatch(d BaseDispatcher) *UploadDispatch {
	return &UploadDispatch{
		BaseDispatcher: d,
	}
}

func (d *UploadDispatch) ProcessTask(tsk *task.Task) error {
	//tsk.SetState(task.StateUpload)
	//d.list.Updated()
	sb, err := d.Sandbox()
	if err != nil {
		return err
	}
	var id string
	if tsk.Type == task.URLTask {
		id, err = sb.SubmitURL(tsk.Path)
	} else {
		id, err = sb.SubmitFile(tsk.Path)
	}
	if err != nil {
		return err
	}
	tsk.SetSandboxID(id)
	logging.Infof("Accepted: %v", id)
	//tsk.SetState(task.StateAccepted)
	tsk.SetChannel(task.ChResult)
	d.list.Updated()
	//d.Channel(ChWait) <- tsk.Number
	return nil
}
