package task

import (
	"encoding/json"
	"examen/pkg/state"
	"testing"
)

func TestList(t *testing.T) {
	l := NewList()
	t0 := NewTask("d:\\abc")
	t0.SetError("ups!")
	l.Add(t0)
	t1 := NewTask("C:\\def")
	t1.SetState(state.StateHighRisk)
	l.Add(t1)
	j, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("JSON: %s", string(j))
}
