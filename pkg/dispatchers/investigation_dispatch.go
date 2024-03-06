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
	//tsk.SetState(task.StateInvestigation)
	//d.list.Updated()
	zipFilePath, err := tsk.InvestigationPath()
	if err != nil {
		return err
	}
	//zipFileName := fmt.Sprintf("%s.zip", tsk.SHA256)
	//zipFilePath := filepath.Join(tasksFolder, zipFileName)
	if err := sbox.GetInvestigation(tsk.SandboxID, zipFilePath); err != nil {
		return err
	}
	tsk.SetInvestigation(zipFilePath)
	tsk.SetChannel(task.ChDone)
	d.list.Updated()
	//taskFileName := "task.json"
	//taskFilePath := filepath.Join(tasksFolder, taskFileName)
	return nil // tsk.SaveToFile(taskFilePath)
}

/*
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
*/
