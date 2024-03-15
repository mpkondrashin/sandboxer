/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

list.go

List of tasks
*/
package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
)

type TaskListInterface interface {
	NewTask(path string) ID
}

type TaskList struct {
	mx         sync.RWMutex
	changed    chan struct{}
	Tasks      map[ID]*Task
	TasksCount ID
}

func NewList() *TaskList {
	l := &TaskList{
		Tasks:   make(map[ID]*Task),
		changed: make(chan struct{}, 1000),
	}
	logging.Debugf("%p XXX List Lock (in New List)", l)
	l.changed <- struct{}{}
	return l
}

func (l *TaskList) Updated() {
	//logging.Debugf("Updated")
	if len(l.changed) > 0 {
		return
	}
	l.changed <- struct{}{}
}

func (l *TaskList) Length() int {
	return len(l.Tasks)
}

func (l *TaskList) Changes() chan struct{} {
	return l.changed
}

var ErrAlreadyExists = errors.New("task already exist")

func (l *TaskList) NewTask(taskType TaskType, path string) (ID, error) {
	defer l.lockUnlock()()
	tsk := l.FindTask(path)
	if tsk != nil {
		tsk.SubmitTime = time.Now()
		l.Updated()
		return 0, fmt.Errorf("%s: %w", path, ErrAlreadyExists)

	}
	logging.Debugf("NewTask %d, %s", l.TasksCount, path)
	tsk = NewTask(l.TasksCount, taskType, path)
	l.Tasks[tsk.Number] = tsk
	l.Updated()
	l.TasksCount++
	return tsk.Number, nil
}

func (l *TaskList) FindTask(path string) *Task {
	// More sophisticated paths comparison algoritm could be used or
	// os.SameFile function
	for _, tsk := range l.Tasks {
		if path == tsk.Path {
			return tsk
		}
	}
	return nil
}

func (l *TaskList) DelByID(id ID) {
	defer l.lockUnlock()() //mx.Lock()
	logging.Debugf("DelByID, id = %d, len = %d", id, len(l.Tasks))
	delete(l.Tasks, id)
	l.Updated()
}

func (l *TaskList) Get(num ID) *Task {
	defer l.lockUnlock()()
	return l.Tasks[num]
}

func (l *TaskList) Task(num ID, callback func(tsk *Task) error) error {
	//defer l.lockUnlock()()
	tsk := l.Tasks[num]
	if tsk == nil {
		return fmt.Errorf("missing task #%d", num)
	}
	//defer tsk.lockUnlock()()
	return callback(tsk)
}

func (l *TaskList) lockUnlock() func() {
	//logging.Debugf("Lock %p", l)
	l.mx.Lock()
	return l.unlock // func() {} //
}
func (l *TaskList) unlock() {
	//logging.Debugf("Unlock %p", l)
	l.mx.Unlock()
}

func (l *TaskList) GetIDs() []ID {
	keys := make([]ID, len(l.Tasks))
	logging.Debugf("keys len = %d", len(l.Tasks))
	i := 0
	for k := range l.Tasks {
		keys[i] = k
		i++
	}
	return keys
}

func (l *TaskList) Process(callback func([]ID)) {
	defer l.lockUnlock()()
	keys := l.GetIDs()
	sort.Slice(keys, func(i, j int) bool {
		return l.Tasks[keys[i]].SubmitTime.After(l.Tasks[keys[j]].SubmitTime)
	})
	//logging.Debugf("slice: %v", keys)
	callback(keys)
}

func (l *TaskList) SortTasks() {
	defer l.lockUnlock()()
	keys := l.GetIDs()
	sort.Slice(keys, func(i, j int) bool {
		return l.Tasks[keys[i]].SubmitTime.After(l.Tasks[keys[j]].SubmitTime)
	})
	//logging.Debugf("slice: %v", keys)
	//l.
}

func (l *TaskList) CountActiveTasks() (count int) {
	for _, t := range l.Tasks {
		if t.Channel != ChDone {
			count++
		}
	}
	return
}

func (l *TaskList) LoadTasks(keepDays int) error {
	folder, err := globals.TasksFolder()
	if err != nil {
		return err
	}
	dir, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	keepDuration := time.Duration(keepDays) * time.Hour * 60
	oldest := time.Now().Add(-keepDuration)
	for _, d := range dir {
		if !d.IsDir() {
			continue
		}
		path := filepath.Join(folder, d.Name(), taskFileName)
		tsk, err := LoadTask(path)
		if err != nil {
			logging.LogError(err)
			continue
		}
		if tsk.SubmitTime.Before(oldest) {
			logging.Debugf("To delete %v", tsk)
			logging.LogError(tsk.Delete())
		}
		tsk.Number = l.TasksCount
		l.TasksCount++
		l.Tasks[tsk.Number] = tsk
	}
	return nil
}

/*
func (l *TaskList) Save(filePath string) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
*/
