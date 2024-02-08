package dispatchers

import (
	"context"
	"fmt"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
	"strings"

	"github.com/mpkondrashin/vone"
)

type ResultDispatch struct {
	BaseDispatcher
}

func NewResultDispatch(d BaseDispatcher) *ResultDispatch {
	return &ResultDispatch{
		BaseDispatcher: d,
	}
}

func (d *ResultDispatch) InboundChannel() int {
	return ChResult
}

func (d *ResultDispatch) ProcessTask(tsk *task.Task) error {
	vOne, err := d.vOne()
	if err != nil {
		return err
	}
	//	time.Sleep(10 * time.Second)
	tsk.SetState(task.StateCheck)
	d.list.Updated()
	results, err := vOne.SandboxAnalysisResults(tsk.SandboxID).Do(context.TODO())
	if err != nil {
		return err
	}
	tsk.SetDigest(results.Digest.MD5, results.Digest.SHA1, results.Digest.SHA256)
	//logging.Debugf("XXX MESSAGE SET: %v", tsk)
	switch results.RiskLevel {
	case vone.RiskLevelHigh:
		tsk.SetRiskLevel(task.RiskLevelHigh)
	case vone.RiskLevelMedium:
		tsk.SetRiskLevel(task.RiskLevelMedium)
	case vone.RiskLevelLow:
		tsk.SetRiskLevel(task.RiskLevelLow)
	case vone.RiskLevelNoRisk:
		tsk.SetRiskLevel(task.RiskLevelNoRisk)
		tsk.SetMessage(task.RiskLevelNoRisk.String())
	default:
		return fmt.Errorf("unknown risk level: %d", results.RiskLevel)
	}
	d.Channel(ChReport) <- tsk.Number
	detectionName := strings.Join(results.DetectionNames, ", ")
	threatType := strings.Join(results.ThreatTypes, ", ")
	tsk.SetMessage(detectionName + threatType)
	logging.Debugf("Type: %s, TrueFileType: %s, RiskLevel: %s, DetectionNames: %s, threatTypes: %s; for task %v",
		results.Type, results.TrueFileType, results.RiskLevel, detectionName, threatType, tsk.Number)
	d.list.Updated()
	return nil
}
