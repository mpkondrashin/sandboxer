package dispatchers

import (
	"context"
	"fmt"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/state"
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
	tsk.SetState(state.StateCheck)
	d.list.Updated()
	results, err := vOne.SandboxAnalysisResults(tsk.SandboxID).Do(context.TODO())
	if err != nil {
		if strings.HasSuffix(err.Error(), "NotFound: Not Found") {
			logging.Debugf("error %v", err)
			d.Channel(ChResult) <- tsk.Number
			//task.SetState(id, state.StateUnsupported)
			return nil
		}
		return err
	}
	//logging.Debugf("XXX MESSAGE SET: %v", tsk)
	switch results.RiskLevel {
	case vone.RiskLevelHigh:
		tsk.SetState(state.StateHighRisk)
	case vone.RiskLevelMedium:
		tsk.SetState(state.StateMediumRisk)
	case vone.RiskLevelLow:
		tsk.SetState(state.StateLowRisk)
	case vone.RiskLevelNoRisk:
		tsk.SetState(state.StateNoRisk)
	default:
		return fmt.Errorf("unknown risk level: %d", results.RiskLevel)
	}
	detectionName := strings.Join(results.DetectionNames, ", ")
	threatType := strings.Join(results.ThreatTypes, ", ")
	tsk.SetMessage(detectionName + threatType)
	logging.Debugf("Type: %s, TrueFileType: %s, RiskLevel: %s, DetectionNames: %s, threatTypes: %s; for task %v",
		results.Type, results.TrueFileType, results.RiskLevel, detectionName, threatType, tsk.Number)
	d.list.Updated()
	return nil
}
