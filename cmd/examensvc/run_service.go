package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"examen/pkg/config"
	"examen/pkg/globals"
	"examen/pkg/logging"
	"examen/pkg/state"
	"examen/pkg/task"

	"github.com/mpkondrashin/vone"
)

type Scan struct {
	list   *task.List
	config *config.Configuration
	vOne   *vone.VOne
	//check  *goperic.Periculosum
}

func NewScan(config *config.Configuration, vOne *vone.VOne /*, check *goperic.Periculosum*/, list *task.List) *Scan {
	return &Scan{
		list:   list,
		config: config,
		vOne:   vOne,
		//check:  check,
	}
}

func (s *Scan) InspecfFolder(folderPath string) {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				s.InspectFile(path)
			}
			return nil
		})
	logging.LogError(err)
}

func (s *Scan) ShouldIgnore(filePath string) bool {
	fileName := filepath.Base(filePath)
	for _, mask := range s.config.Ignore {
		result, err := filepath.Match(strings.ToLower(mask), strings.ToLower(fileName))
		logging.LogError(err)
		if result {
			logging.Debugf("%s: ignore by mask \"%s\"", filePath, mask)
			return true
		}
	}
	return false
}

func (s *Scan) InspectFile(filePath string) {
	logging.Debugf("InspectFile(%s)", filePath)
	info, err := os.Lstat(filePath)
	if err != nil {
		logging.Errorf("%s", err)
	}
	if info.IsDir() {
		s.InspecfFolder(filePath)
		return
	}
	if !info.Mode().IsRegular() {
		logging.Errorf("%s: not regular file", filePath)
	}
	if s.ShouldIgnore(filePath) {
		return
	}
	t := task.NewTask(filePath)
	s.list.Add(t)
	if err := s.Submit(t); err != nil {
		logging.LogError(err)
		t.SetError(err.Error())
		return
	}
}

func (s *Scan) Submit(t *task.Task) error {
	t.SetState(state.StateUpload)
	f, err := s.vOne.SandboxSubmitFile().SetFilePath(t.Path)
	if err != nil {
		return err
	}
	response, headers, err := f.Do(context.TODO())
	_ = headers
	if err != nil {
		return err
	}
	t.SetID(response.ID)
	logging.Infof("Accepted: %v", t)
	if err := s.WaitForResult(t); err != nil {
		return err
	}
	//c.LogQuota(id, headers)
	if err := s.GetResult(t); err != nil {
		return err
	}
	return nil
}

func (s *Scan) WaitForResult(t *task.Task) error {
	t.SetState(state.StateInspect)
	for {
		// Should we set temporary state "checking"?
		status, err := s.vOne.SandboxSubmissionStatus(t.VOneID()).Do(context.TODO())
		if err != nil {
			return fmt.Errorf("check status: %w", err)
		}
		logging.Debugf("%s Status: %v", t.VOneID(), status.Status)
		switch status.Status {
		case vone.StatusSucceeded:
			return nil
		case vone.StatusRunning:
			//if time.Now().After(endTime) {
			//	return ErrTimeout
			//}
			time.Sleep(s.config.Sleep)
		case vone.StatusFailed:
			return fmt.Errorf("%s: %s", status.Error.Code, status.Error.Message)
		default:
			return fmt.Errorf("unknown status: %s", status.Status)
		}
	}
}

func (s *Scan) GetResult(t *task.Task) error {
	results, err := s.vOne.SandboxAnalysisResults(t.VOneID()).Do(context.TODO())
	if err != nil {
		return err
	}
	detectionName := strings.Join(results.DetectionNames, ", ")
	threatType := strings.Join(results.ThreatTypes, ", ")
	logging.Debugf("Type: %s, TrueFileType: %s, RiskLevel: %s, DetectionNames: %s, threatTypes: %s; for task %v",
		results.Type, results.TrueFileType, results.RiskLevel, detectionName, threatType, t)
	switch results.RiskLevel {
	case vone.RiskLevelHigh:
		t.SetState(state.StateHighRisk)
	case vone.RiskLevelMedium:
		t.SetState(state.StateMediumRisk)
	case vone.RiskLevelLow:
		t.SetState(state.StateLowRisk)
	case vone.RiskLevelNoRisk:
		t.SetState(state.StateNoRisk)
	default:
		err := fmt.Errorf("unknown risk level: %d", results.RiskLevel)
		logging.LogError(err)
		t.SetError(err.Error())
	}
	return nil
}

const examenSvcLog = "examen_svc.log"

func RunService() (func(), error) {
	conf, err := config.LoadConfiguration(globals.AppID, globals.ConfigFileName)
	if err != nil {
		return nil, err
	}
	/*close := logging.NewFileLog(conf.LogFolder(), examenSvcLog)
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	*/
	inbox := make(StringChannel)
	list := task.NewList()
	go SubmitDispatch(inbox)
	//	pericPath, err := config.PericulosumPath()

	//	if err != nil {
	//	}
	//	goperic.NewPericulosum()
	vOne := vone.NewVOne(conf.Domain, conf.Token)
	scan := NewScan(conf, vOne, list)
	go func() {
		for {
			s := <-inbox
			logging.Debugf("Got %s", s)
			go scan.InspectFile(s)
		}
	}()
	return nil, nil
}
