/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

report_dispatch.go

Download PDF report dispatcher
*/
package dispatchers

import (
	"sandboxer/pkg/task"
)

type ReportDispatch struct {
	BaseDispatcher
}

func NewReportDispatch(d BaseDispatcher) *ReportDispatch {
	return &ReportDispatch{
		BaseDispatcher: d,
	}
}

func (d *ReportDispatch) InboundChannel() task.Channel {
	return task.ChReport
}

func (d *ReportDispatch) ProcessTask(tsk *task.Task) error {
	sbox, err := d.Sandbox()
	if err != nil {
		return err
	}
	filePath, err := tsk.ReportPath()
	if err != nil {
		return err
	}
	if err := sbox.GetReport(tsk.SandboxID, filePath); err != nil {
		return err
	}
	tsk.SetReport(filePath)
	d.list.Updated()
	tsk.SetChannel(task.ChInvestigation)
	return nil
}

/*
func (d *ReportDispatch) ReportPath(tsk *task.Task) (string, error) {
	baseFolder, err := globals.UserDataFolder()
	if err != nil {
		return "", err
	}
	folder := filepath.Join(baseFolder, globals.TasksFolder, tsk.SHA256)
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s.pdf", tsk.SHA256)
	return filepath.Join(folder, fileName), nil
}
*/
