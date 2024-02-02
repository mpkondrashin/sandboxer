package state

//go:generate enum -package state -type State -names New,Upload,Inspected,Check,WaitForResult,Ignored,Unsupported,Error,NoRisk,LowRisk,MediumRisk,HighRisk,Count

//var ErrIDNotFound = errors.New("id not found")

/*
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
*/
