/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

upload_dispatch.go

Upload file for inspection
*/
package dispatchers

import (
	"context"
	"fmt"
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
	vOne, err := d.vOne()
	if err != nil {
		return err
	}
	if tsk.Type == task.URLTask {
		f := vOne.SandboxSubmitURLs().AddURL(tsk.Path)
		response, _, err := f.Do(context.TODO())
		if err != nil {
			return err
		}
		if len(response) != 1 {
			return fmt.Errorf("wrong response length: %v", response)
		}
		tsk.SetSandboxID(response[0].Body.ID)
		logging.Infof("Accepted: %v", response[0].Body.ID)
	} else {
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
	}
	//tsk.SetState(task.StateAccepted)
	tsk.SetChannel(task.ChWait)
	d.list.Updated()
	//d.Channel(ChWait) <- tsk.Number
	return nil
}
