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

//type TasksMap Map[ID, *Task]

type TaskList struct {
	mx         sync.RWMutex
	changed    chan struct{}
	Tasks      *Map[ID, *Task]
	TasksCount ID
}

func NewList() *TaskList {
	l := &TaskList{
		Tasks:   new(Map[ID, *Task]),
		changed: make(chan struct{}, 1000),
	}
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
	return l.Tasks.Length()
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
	//l.Tasks[tsk.Number] = tsk
	l.Tasks.Store(tsk.Number, tsk)
	l.Updated()
	l.TasksCount++
	return tsk.Number, nil
}

func (l *TaskList) FindTask(path string) (result *Task) {
	// More sophisticated paths comparison algoritm could be used or
	// os.SameFile function
	l.Tasks.Range(func(id ID, tsk *Task) bool {
		if path == tsk.Path {
			result = tsk
			return false
		}
		return true
	})
	return
}

func (l *TaskList) DelByID(id ID) {
	//defer l.lockUnlock()() //mx.Lock()
	//logging.Debugf("DelByID, id = %d, len = %d", id, len(l.Tasks))
	l.Tasks.Delete(id)
	l.Updated()
}

func (l *TaskList) Get(num ID) *Task {
	//
	// It was here:
	//defer l.lockUnlock()()
	tsk, _ := l.Tasks.Load(num)
	return tsk
}

func (l *TaskList) Task(num ID, callback func(tsk *Task) error) error {
	//defer l.lockUnlock()()
	tsk, ok := l.Tasks.Load(num)
	if !ok {
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

func (l *TaskList) GetIDs() (keys []ID) {
	l.Tasks.Range(func(id ID, _ *Task) bool {
		keys = append(keys, id)
		return true
	})
	return
}

func (l *TaskList) Process(callback func([]ID)) {
	//defer l.lockUnlock()()
	keys := l.GetIDs()
	sort.Slice(keys, func(i, j int) bool {
		a, ok := l.Tasks.Load(keys[i])
		if !ok {
			return false
		}
		b, ok := l.Tasks.Load(keys[j])
		if !ok {
			return false
		}
		return a.SubmitTime.After(b.SubmitTime)
	})
	callback(keys)
}

/*
	func (l *TaskList) SortTasks() {
		defer l.lockUnlock()()
		keys := l.GetIDs()
		sort.Slice(keys, func(i, j int) bool {
			a, ok := l.Tasks.Load(keys[i])
			if !ok {
				return false
			}
			b, ok := l.Tasks.Load(keys[j])
			if !ok {
				return false
			}
			return a.SubmitTime.After(b.SubmitTime)
		})
	}
*/

func (l *TaskList) CountActiveTasks() (count int) {
	l.Tasks.Range(func(id ID, tsk *Task) bool {
		if tsk.Channel != ChDone {
			count++
		}
		return true
	})
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
		l.Tasks.Store(tsk.Number, tsk)
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

func (l *TaskList) DeleteTask(tsk *Task) error {
	err := tsk.Delete()
	if err != nil {
		return err
	}
	l.DelByID(tsk.Number)
	return nil
}

func (l *TaskList) DeleteSameTasks(tsk *Task) (err error) {
	l.Tasks.Range(func(id ID, t *Task) bool {
		if t.Channel != tsk.Channel {
			return true
		}
		if t.RiskLevel != tsk.RiskLevel {
			return true
		}
		err = t.Delete()
		if err != nil {
			return false
		}
		l.DelByID(t.Number)
		return true
	})
	return
}

func (l *TaskList) DeleteAllTasks() (err error) {
	l.Tasks.Range(func(id ID, t *Task) bool {
		err = t.Delete()
		if err != nil {
			return false
		}
		l.DelByID(t.Number)
		return true
	})
	return
}
