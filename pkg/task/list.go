package task

import (
	"examen/pkg/logging"
	"examen/pkg/state"
	"sort"
	"sync"
)

var _list *TaskList

func init() {
	_list = NewList()
}

func List() *TaskList {
	return _list
}

func SetState(id ID, s state.State) {
	_list.Get(id).SetState(s)
}

func SetSandboxID(id ID, vOneID string) {
	_list.Get(id).SetID(vOneID)
}

func GetSandboxID(id ID) string {
	return _list.Get(id).SandboxID
}

func SetError(id ID, err error) {
	_list.Get(id).SetError(err)
}

func Path(id ID) string {
	return _list.Get(id).Path
}

func Delete(id ID) {
	_list.DelByID(id)
}

func New(path string) ID {
	t := NewTask(path)
	_list.Add(t)
	return t.Number
}

/*func Add(t *Task) {
	_list.Add(t)
}
func Del(t *Task) {
	_list.Del(t)
}

func Iterate(callback func(*Task)) {
	_list.Iterate(callback)
}*/

type TaskList struct {
	mx sync.RWMutex
	//changeMX sync.Mutex
	changed chan struct{}
	Tasks   map[ID]*Task
}

func NewList() *TaskList {
	l := &TaskList{
		Tasks:   make(map[ID]*Task),
		changed: make(chan struct{}, 1000),
	}
	logging.Debugf("%p XXX List Lock (in New List)", l)
	//l.changeMX.Lock()
	l.changed <- struct{}{}
	return l
}

func (l *TaskList) Updated() {
	if len(l.changed) > 0 {
		return
	}
	l.changed <- struct{}{}
}

func (l *TaskList) Changes() chan struct{} {
	return l.changed
}

func (l *TaskList) Add(tsk *Task) {
	l.mx.Lock()
	defer l.mx.Unlock()
	l.Tasks[tsk.Number] = tsk
	logging.Debugf("%p XXX List Unlock (in Add)", l)
	l.Updated()
}

func (l *TaskList) Del(tsk *Task) {
	l.mx.Lock()
	defer l.mx.Unlock()
	delete(l.Tasks, tsk.Number)
	l.Updated()
}
func (l *TaskList) DelByID(id ID) {
	l.mx.Lock()
	defer l.mx.Unlock()
	delete(l.Tasks, id)
	l.Updated()
}

func (l *TaskList) Get(num ID) *Task {
	l.mx.RLock()
	defer l.mx.RUnlock()
	return l.Tasks[num]
}

func (l *TaskList) Iterate(callback func(*Task)) {
	l.mx.RLock()
	defer l.mx.RUnlock()
	keys := make([]ID, len(l.Tasks))
	i := 0
	for k := range l.Tasks {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		callback(l.Tasks[k])
	}
}

func (l *TaskList) IterateIDs(from int, count int, callback func(id ID)) {
	l.mx.RLock()
	defer l.mx.RUnlock()
	keys := make([]ID, len(l.Tasks))
	i := 0
	for k := range l.Tasks {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	if from >= len(keys) {
		from = len(keys) - count
	}
	if from < 0 {
		from = 0
	}
	if from+count > len(keys) {
		count = len(keys) - from
	}
	logging.Debugf("from: %d, count: %d", from, count)
	for _, k := range keys { //[from : from+count] {
		callback(k)
	}
}

/*
func (l *TaskList) WaitForChange() {
	logging.Debugf("%p XXX List Lock (in Wait)", l)
	l.Updated()
}
*/
