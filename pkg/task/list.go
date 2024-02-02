package task

import (
	"fmt"
	"sort"
	"sync"

	"sandboxer/pkg/logging"
)

/*
var _list *TaskList

func init() {
	_list = NewList()
}

func List() *TaskList {
	return _list
}

func SetState(id ID, s state.State) {
	_list.Get(id).SetState(s)
	_list.Updated()
}

func SetSandboxID(id ID, vOneID string) {
	_list.Get(id).SetID(vOneID)
}

func GetSandboxID(id ID) string {
	return _list.Get(id).SandboxID
}

func SetError(id ID, err error) {
	_list.Get(id).SetError(err)
	_list.Updated()
}
func SetMessage(id ID, msg string) {
	_list.Get(id).SetMessage(msg)
	_list.Updated()
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

/*
Submissions:
task.Delete(tsk.Number)
task.List().Process(func(ids []task.ID) {
return s.CardWidget(task.NewTask("placeholder"))
task.List().Task(ids[i], func(tsk *task.Task) error {
tsk = task.NewTask("placeholder")
<-task.List().Changes():

Dispatchers:
task.SetMessage(id, detectionName+threatType)
task.New(s)
task.Path(id)
task.SetSandboxID(id, response.ID)
task.Delete(id)
task.GetSandboxID(id)
task.SetError(id, err)
*/
type TaskListInterface interface {
	NewTask(path string) ID
}

type TaskList struct {
	mx sync.RWMutex
	//changeMX sync.Mutex
	changed    chan struct{}
	Tasks      map[ID]*Task
	tasksCount ID
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

func (l *TaskList) NewTask(path string) ID {
	defer l.lockUnlock()() //mx.Lock()
	logging.Debugf("NewTask %d, %s", l.tasksCount, path)
	tsk := NewTask(l.tasksCount, path)
	l.Tasks[tsk.Number] = tsk
	l.Updated()
	l.tasksCount++
	return tsk.Number
}

func (l *TaskList) Del(tsk *Task) {
	defer l.lockUnlock()()
	delete(l.Tasks, tsk.Number)
	l.Updated()
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

func (l *TaskList) Iterate(callback func(*Task)) {
	defer l.lockUnlock()()
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

/*
func (l *TaskList) IterateIDs(from int, count int, callback func(id ID)) {
	defer l.lockUnlock()()

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
*/

func (l *TaskList) lockUnlock() func() {
	logging.Debugf("Lock %p", l)
	l.mx.Lock()
	return l.unlock // func() {} //
}
func (l *TaskList) unlock() {
	logging.Debugf("Unlock %p", l)
	l.mx.Unlock()
}

func (l *TaskList) GetIDs() []ID {
	defer l.lockUnlock()()
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
	keys := l.GetIDs()
	sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	logging.Debugf("slice: %v", keys)
	callback(keys)
}

/*
func (l *TaskList) WaitForChange() {
	logging.Debugf("%p XXX List Lock (in Wait)", l)
	l.Updated()
}
*/
