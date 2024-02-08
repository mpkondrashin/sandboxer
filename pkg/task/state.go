package task

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type State int

const (
	StateNew State = iota
	StateUpload
	StateInspected
	StateCheck
	StateWaitForResult
	StateReport
	StateInvestigation
	StateDone
	StateCount
)

// String - return string representation for State value
func (v State) String() string {
	s, ok := map[State]string{
		StateNew:           "New",
		StateUpload:        "Upload",
		StateInspected:     "Inspected",
		StateCheck:         "Check",
		StateWaitForResult: "Wait For Result",
		StateReport:        "Report",
		StateInvestigation: "Investiation",
		StateDone:          "Done",
		StateCount:         "Count",
	}[v]
	if ok {
		return s
	}
	return "State(" + strconv.FormatInt(int64(v), 10) + ")"
}

// ErrUnknownState - will be returned wrapped when parsing string
// containing unrecognized value.
var ErrUnknownState = errors.New("unknown State")

var mapStateFromString = map[string]State{
	"new":             StateNew,
	"upload":          StateUpload,
	"inspected":       StateInspected,
	"check":           StateCheck,
	"wait for result": StateWaitForResult,
	"report":          StateReport,
	"investigation":   StateInvestigation,
	"finished":        StateDone,
	"count":           StateCount,
}

// UnmarshalJSON implements the Unmarshaler interface of the json package for State.
func (s *State) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	result, ok := mapStateFromString[strings.ToLower(v)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownState, v)
	}
	*s = result
	return nil
}

// MarshalJSON implements the Marshaler interface of the json package for State.
func (s State) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", s)), nil
}

// UnmarshalYAML implements the Unmarshaler interface of the yaml.v3 package for State.
func (s *State) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	result, ok := mapStateFromString[strings.ToLower(v)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownState, v)
	}
	*s = result
	return nil
}

// MarshalXML implements the Marshaler interface of the xml package for State.
func (s State) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(s.String(), start)
}

// UnmarshalXML implements the Unmarshaler interface of the xml package for State.
func (s *State) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	result, ok := mapStateFromString[strings.ToLower(v)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownState, v)
	}
	*s = result
	return nil
}
