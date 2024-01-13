package state

import "errors"

//go:generate enum -package state -type State -names Unknown,Upload,Inspect,Report,Unsupported,Error,NoRisk,LowRisk,MediumRisk,HighRisk,Count

var ErrIDNotFound = errors.New("id not found")

type Object struct {
	ID    string
	Path  string
	State State
}

func NewObject(id, path string) Object {
	return Object{ID: id, Path: path, State: StateUnknown}
}

type FileState interface {
	AddObject(Object) error
	SetState(id string, state State) error
	ListObjects() ([]Object, error)
}
