/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

upload_dispatch.go

Upload file for inspection
*/
package dispatchers

import (
	"context"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
)

type UploadDispatch struct {
	BaseDispatcher
}

func (*UploadDispatch) InboundChannel() int {
	return ChUpload
}

func NewUploadDispatch(d BaseDispatcher) *UploadDispatch {
	return &UploadDispatch{
		BaseDispatcher: d,
	}
}

func (d *UploadDispatch) ProcessTask(tsk *task.Task) error {
	tsk.SetState(task.StateUpload)
	d.list.Updated()

	vOne, err := d.vOne()
	if err != nil {
		return err
	}
	f, err := vOne.SandboxSubmitFile().SetFilePath(tsk.Path)
	if err != nil {
		return err
	}
	response, _, err := f.Do(context.TODO())
	if err != nil {
		return err
	}
	tsk.SetSandboxID(response.ID)
	logging.Infof("Accepted: %v", response.ID)
	tsk.SetState(task.StateInspected)
	d.list.Updated()
	d.Channel(ChWait) <- tsk.Number
	return nil
}
