/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

submit_dispatch.go

Send files from FIFO to prefilter channel
*/
package dispatchers

import (
	"log"
	"sync"

	"sandboxer/pkg/fifo"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
)

const StopPath = "STOP"

type SubmitDispatch struct {
	BaseDispatcher
}

func NewSubmitDispatch(b BaseDispatcher) *SubmitDispatch {
	return &SubmitDispatch{b}
}

func (d *SubmitDispatch) Run(wg *sync.WaitGroup) {
	logging.Debugf("Start %T", d)
	for {
		fifoReader, err := fifo.NewReader()
		if err != nil {
			log.Fatal(err)
		}
		var s string
		if err := fifoReader.Read(&s); err != nil {
			if err.Error() != "EOF" {
				logging.Errorf("read FIFO: %v", err)
			}
			continue
		}
		fifoReader.Close()
		logging.Infof("Got new path: %s", s)
		if s == StopPath {
			break
		}
		//		logging.Debugf("SEND %s to %d", s, ChPrefilter)
		tsk, err := d.list.NewTask(task.FileTask, s)
		if err != nil {
			//if !errors.Is(err, task.ErrAlreadyExists) {
			logging.LogError(err)
			//}
			continue
		}
		d.Channel(task.ChPrefilter) <- tsk
		//d.list.Updated()
	}
	wg.Done()
	logging.Debugf("Stop SubmitDispatch")
}

/*
type StringChannel chan string

func SubmitDispatchFunc(ch StringChannel) {
	for {
		fifoReader, err := fifo.NewReader()
		if err != nil {
			log.Fatal(err)
		}
		var s string
		if err := fifoReader.Read(&s); err != nil {
			if err.Error() != "EOF" {
				logging.Errorf("read FIFO: %v", err)
			}
			continue
		}
		fifoReader.Close()
		logging.Infof("Got new path: %s", s)
		ch <- s
	}
}
*/
