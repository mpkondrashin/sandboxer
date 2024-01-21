package grpc

import (
	"examen/pkg/task"
	"testing"
)

func TestClientServer(t *testing.T) {
	path := "C:\\a\file.exe"
	t1 := task.NewTask("C:\\a\\b.txt")
	t2 := task.NewTask("C:\\c\\d.txt")
	go func() {
		err := RunServer(
			func(name string) error {
				if name != path {
					t.Errorf("%s != %s", name, path)
				}
				return nil
			},
			func(from, count int32) ([]*task.Task, error) {
				return []*task.Task{t1, t2}, nil
			},
		)
		if err != nil {
			t.Log(err)
		}
	}()
	err := Submit(path)
	if err != nil {
		t.Error(err)
	}
	i := 0
	err = Status(0, 1, func(tsk *task.Task) {
		if i == 0 {
			if tsk.Path != t1.Path {
				t.Errorf("%s != %s", tsk.Path, t1.Path)
			}
			i++
		} else {
			if tsk.Path != t2.Path {
				t.Errorf("%s != %s", tsk.Path, t2.Path)
			}
		}
		//t.Log(tsk)
	})
	if err != nil {
		t.Error(err)
	}
}
