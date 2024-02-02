package dispatchers

import (
	"errors"
	"os"
	"path/filepath"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/state"
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

func (*PrefilterDispatch) InboundChannel() int {
	return ChPrefilter
}

func (d *PrefilterDispatch) ProcessTask(tsk *task.Task) error {
	logging.Debugf("Prefilter %s", tsk.Path)
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
	if d.ShouldIgnore(tsk.Path) {
		tsk.SetState(state.StateIgnored)
		d.list.Updated()
		return nil
	}
	logging.Debugf("Send Task #%d to %d", tsk.Number, ChUpload)
	d.Channel(ChUpload) <- tsk.Number
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
			p.Channel(ChPrefilter) <- p.list.NewTask(path)
			return nil
		})
	logging.LogError(err)
	// XXX id := task.New(folderPath)
	// XXX task.SetError(id, err)
}

func (p *PrefilterDispatch) ShouldIgnore(filePath string) bool {
	fileName := filepath.Base(filePath)
	for _, mask := range p.conf.Ignore {
		result, err := filepath.Match(strings.ToLower(mask), strings.ToLower(fileName))
		logging.LogError(err)
		if result {
			logging.Debugf("%s: ignore by mask \"%s\"", filePath, mask)
			return true
		}
	}
	return false
}
