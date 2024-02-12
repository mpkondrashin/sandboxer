/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

investigation_dispatch.go

Download investigation package
*/
package dispatchers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sandboxer/pkg/globals"
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

func (d *InvestigationDispatch) InboundChannel() int {
	return ChInvestigation
}

func (d *InvestigationDispatch) ProcessTask(tsk *task.Task) error {
	vOne, err := d.vOne()
	if err != nil {
		return err
	}
	tsk.SetState(task.StateInspected)
	d.list.Updated()
	tasksFolder, err := d.TasksPath(tsk)
	if err != nil {
		return err
	}
	zipFileName := fmt.Sprintf("%s.zip", tsk.SHA256)
	zipFilePath := filepath.Join(tasksFolder, zipFileName)
	if err := vOne.SandboxInvestigationPackage(tsk.SandboxID).Store(context.TODO(), zipFilePath); err != nil {
		return err
	}
	tsk.SetInvestigation(zipFilePath)
	tsk.SetState(task.StateDone)
	d.list.Updated()
	taskFileName := "task.yaml"
	taskFilePath := filepath.Join(tasksFolder, taskFileName)
	return tsk.Save(taskFilePath)
}

func (d *InvestigationDispatch) TasksPath(tsk *task.Task) (string, error) {
	baseFolder, err := globals.UserDataFolder()
	if err != nil {
		return "", err
	}
	folder := filepath.Join(baseFolder, globals.TasksFolder, tsk.SHA256)
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}
	return folder, nil
}
