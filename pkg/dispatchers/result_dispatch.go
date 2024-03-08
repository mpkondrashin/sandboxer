/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

result_dispatch.go

Get inspection result
*/
package dispatchers

import (
	"fmt"
	"path/filepath"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/sandbox"
	"sandboxer/pkg/task"
	"sandboxer/pkg/xplatform"
	"time"
)

type ResultDispatch struct {
	BaseDispatcher
}

func NewResultDispatch(d BaseDispatcher) *ResultDispatch {
	return &ResultDispatch{
		BaseDispatcher: d,
	}
}

func (d *ResultDispatch) InboundChannel() task.Channel {
	return task.ChResult
}

func (d *ResultDispatch) ProcessTask(tsk *task.Task) error {
	sb, err := d.Sandbox()
	if err != nil {
		return err
	}
	riskLevel, threatName, err := sb.GetResult(tsk.SandboxID)
	logging.Debugf("GetResut: %v (%d), %s [%v]", riskLevel, riskLevel, threatName, err)
	tsk.SetRiskLevel(riskLevel)
	switch riskLevel {
	case sandbox.RiskLevelNotReady:
		tsk.Deactivate()
		d.list.Updated()
		logging.Debugf("Seleep %v for %v", d.conf.Sleep, tsk)
		time.Sleep(d.conf.Sleep)
		tsk.SetChannel(task.ChResult)
	case sandbox.RiskLevelUnsupported:
		if err != nil {
			tsk.SetMessage(err.Error())
		} else {
			tsk.SetMessage("Unsupported file type")
		}
		tsk.SetChannel(task.ChDone)
		return nil
		//	case sandbox.RiskLevelError:
		//		return err
	default:
		tsk.SetMessage(threatName)
		tsk.SetChannel(task.ChReport)
		if tsk.RiskLevel != sandbox.RiskLevelNoRisk {
			subtitle := fmt.Sprintf("%v threat found %s", tsk.RiskLevel, threatName)
			xplatform.Alert(globals.AppID, globals.AppName, subtitle, filepath.Base(tsk.Path))
		}
	}
	return err
}
