/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

launcher.go

Run all dispatchers
*/
package dispatchers

import (
	"sandboxer/pkg/config"
	"sandboxer/pkg/fifo"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
	"sync"
)

const (
	PrefilterDispatchers = 5
	UploadDispatchers    = 5
	//	WaitDispatchers         = 5
	ResultDispatchers       = 5
	ReportDispatcher        = 5
	InvestigationDispatcher = 5
)

type Launcher struct {
	conf     *config.Configuration
	channels *task.Channels
	list     *task.TaskList
}

func NewLauncher(conf *config.Configuration, channels *task.Channels, list *task.TaskList) *Launcher {
	return &Launcher{
		conf:     conf,
		channels: channels,
		list:     list,
	}
}

func (l *Launcher) Run() {
	base := NewBaseDispatcher(l.conf, l.channels, l.list)
	dispatchers := []struct {
		count      int
		dispatcher Dispatcher
	}{
		{InvestigationDispatcher, NewInvestigationDispatch(base)},
		{ReportDispatcher, NewReportDispatch(base)},
		{ResultDispatchers, NewResultDispatch(base)},
		//{WaitDispatchers, NewWaitDispatch(base)},
		{UploadDispatchers, NewUploadDispatch(base)},
		{PrefilterDispatchers, NewPrefilterDispatch(base)},
	}
	var wg sync.WaitGroup
	for _, d := range dispatchers {
		for i := 0; i < d.count; i++ {
			wg.Add(1)
			go l.RunDispatcher(d.dispatcher, &wg)
		}
	}
	submit := NewSubmitDispatch(base)
	wg.Add(1)
	go submit.Run(&wg)
	l.LoadTasks()
	//wg.Wait()
}

func (l *Launcher) LoadTasks() {
	logging.Debugf("LoadTasks")
	err := l.list.LoadTasks(l.conf.GetTasksKeepDays())
	if err != nil {
		logging.LogError(err)
		return
	}
	logging.Debugf("LoadTasks %d", len(l.list.Tasks))
	l.list.Process(func(ids []task.ID) {
		for _, id := range ids {
			err := l.list.Task(id, func(tsk *task.Task) error {
				logging.Debugf("Process task: %v", tsk)
				if tsk.Channel == task.ChDone {
					return nil
				}
				l.channels.TaskChannel[tsk.Channel] <- tsk.Number
				return nil
			})
			logging.LogError(err)
		}
	})

}

func (l *Launcher) RunDispatcher(disp Dispatcher, wg *sync.WaitGroup) {
	//ctx, cancel := context.WithCancel(context.TODO())
	logging.Debugf("Start %T", disp)
	ch := disp.InboundChannel()
	for id := range l.channels.TaskChannel[ch] {
		_ = l.list.Task(id, func(tsk *task.Task) error {
			logging.Debugf("Got from %v task %v", ch, tsk)
			tsk.Activate()
			logging.Debugf("Activate")
			l.list.Updated()
			err := disp.ProcessTask(tsk)
			logging.Debugf("Deactivate")
			tsk.Deactivate()
			l.list.Updated()
			if err != nil {
				tsk.SetError(err)
				//l.list.Updated()
				logging.Errorf("Task #%d: %v (%T)", id, err, disp)
				return nil
			}
			if tsk.Channel == task.ChDone {
				return nil
			}
			l.channels.TaskChannel[tsk.Channel] <- tsk.Number
			return nil
		})
	}
	wg.Done()
}

func (l *Launcher) Stop() error {
	l.channels.Close() // Should we move it to the end?
	fifoWriter, err := fifo.NewWriter()
	if err != nil {
		return err
	}
	defer func() {
		logging.LogError(fifoWriter.Close())
	}()
	if err = fifoWriter.Write(StopPath); err != nil {
		return err
	}
	return nil
}
