/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

prefilter_dispatch.go

Prefilter tasks
*/
package dispatchers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/sandbox"
	"sandboxer/pkg/task"
	"strings"
)

type PrefilterDispatch struct {
	BaseDispatcher
}

func NewPrefilterDispatch(d BaseDispatcher) *PrefilterDispatch {
	return &PrefilterDispatch{
		BaseDispatcher: d,
	}
}

func (*PrefilterDispatch) InboundChannel() task.Channel {
	return task.ChPrefilter
}

func (d *PrefilterDispatch) ProcessTask(tsk *task.Task) error {
	logging.Debugf("Prefilter %s", tsk.Path)
	if tsk.Type == task.FileTask {
		info, err := os.Lstat(tsk.Path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			d.list.DelByID(tsk.Number)
			go d.InspecfFolder(tsk.Path)
			d.list.Updated()
			return nil
		}
		if !info.Mode().IsRegular() {
			return errors.New("not regular file")
		}
		if err := tsk.CalculateHash(); err != nil {
			return err
		}
		mask := d.MatchIgnoreMask(tsk.Path)
		if mask != "" {
			tsk.SetChannel(task.ChDone)
			tsk.SetRiskLevel(sandbox.RiskLevelUnsupported)
			tsk.SetMessage(fmt.Sprintf("Matched ignore mask '%s'", mask))
			tsk.SetChannel(task.ChDone)
			d.list.Updated()
			return nil
		}
	} else {
		if err := tsk.CalculateHash(); err != nil {
			return err
		}
	} //	logging.Debugf("Send Task #%d to %d", tsk.Number, ChUpload)
	tsk.SetChannel(task.ChSubmit)
	return nil
}

func (p *PrefilterDispatch) InspecfFolder(folderPath string) {
	logging.Debugf("InspectFolder(%s)", folderPath)
	err := filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if filepath.Base(path) == ".git" {
					return filepath.SkipDir
				}
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			tsk, err := p.list.NewTask(task.FileTask, path)
			if err != nil {
				if errors.Is(err, task.ErrAlreadyExists) {
					return nil
				}
				return err
			}
			p.Channel(task.ChPrefilter) <- tsk
			return nil
		})
	logging.LogError(err)
	// XXX id := task.New(folderPath)
	// XXX task.SetError(id, err)
}

func (p *PrefilterDispatch) MatchIgnoreMask(filePath string) string {
	fileName := filepath.Base(filePath)
	for _, mask := range p.conf.Ignore {
		result, err := filepath.Match(strings.ToLower(mask), strings.ToLower(fileName))
		logging.LogError(err)
		if result {
			logging.Debugf("%s: ignore by mask \"%s\"", filePath, mask)
			return mask
		}
	}
	return ""
}
