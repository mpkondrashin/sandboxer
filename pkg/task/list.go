package task

import (
	"sort"
	"sync"
)

type List struct {
	mx    sync.Mutex
	Tasks map[int64]*Task
}

func NewList() *List {
	return &List{
		Tasks: make(map[int64]*Task),
	}
}

func (l *List) Add(task *Task) {
	l.mx.Lock()
	defer l.mx.Unlock()
	l.Tasks[task.Number] = task
}

func (l *List) Del(task *Task) {
	l.mx.Lock()
	defer l.mx.Unlock()
	delete(l.Tasks, task.Number)
}

func (l *List) Iterate(callback func(*Task)) {
	l.mx.Lock()
	defer l.mx.Unlock()
	keys := make([]int64, len(l.Tasks))
	i := 0
	for k, _ := range l.Tasks {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		callback(l.Tasks[k])
	}
}
