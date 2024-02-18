/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

report_dispatch.go

Download PDF report dispatcher
*/
package dispatchers

import (
	"context"
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

func (d *ReportDispatch) InboundChannel() int {
	return ChReport
}

func (d *ReportDispatch) ProcessTask(tsk *task.Task) error {
	vOne, err := d.vOne()
	if err != nil {
		return err
	}
	tsk.SetState(task.StateReport) // Duplicated in result_dispatch
	d.list.Updated()
	filePath, err := tsk.ReportPath()
	if err != nil {
		return err
	}
	if err := vOne.SandboxDownloadResults(tsk.SandboxID).Store(context.TODO(), filePath); err != nil {
		return err
	}
	tsk.SetReport(filePath)
	tsk.SetState(task.StateInvestigation)
	d.list.Updated()
	d.Channel(ChInvestigation) <- tsk.Number
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
