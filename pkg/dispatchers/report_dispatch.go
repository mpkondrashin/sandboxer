package dispatchers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sandboxer/pkg/globals"
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
	tsk.SetState(task.StateReport)
	d.list.Updated()
	filePath, err := d.ReportPath(tsk)
	if err != nil {
		return err
	}
	if err := vOne.SandboxDownloadResults(tsk.SandboxID).Store(context.TODO(), filePath); err != nil {
		return err
	}
	tsk.SetReport(filePath)
	tsk.SetState(task.StateDone)
	d.list.Updated()
	return nil
}

const ReportsFolder = "reports"

func (d *ReportDispatch) ReportPath(tsk *task.Task) (string, error) {
	baseFolder, err := globals.UserDataFolder()
	if err != nil {
		return "", err
	}
	folder := filepath.Join(baseFolder, ReportsFolder, tsk.SHA256)
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s.pdf", tsk.SHA256)
	return filepath.Join(folder, fileName), nil
}
