package state

import "fmt"

type MemState struct {
	objects []Object
}

func NewMemState() *MemState {
	return &MemState{}
}

func (m *MemState) AddObject(o Object) error {
	m.objects = append(m.objects, o)
	return nil
}

func (m *MemState) SetState(id string, state State) error {
	for i := range m.objects {
		if m.objects[i].ID == id {
			m.objects[i].State = state
			return nil
		}
	}
	return fmt.Errorf("SetState: %s: %w", id, ErrIDNotFound)
}

func (m *MemState) ListObjects() ([]Object, error) {
	return m.objects, nil
}
