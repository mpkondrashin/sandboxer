/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

type.go

Types of inspection tasks
*/
package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type TaskType int

const (
	FileTask TaskType = iota
	URLTask
)

// String - return string representation for State value
func (t TaskType) String() string {
	s, ok := map[TaskType]string{
		FileTask: "File",
		URLTask:  "URL",
	}[t]
	if ok {
		return s
	}
	return "State(" + strconv.FormatInt(int64(t), 10) + ")"
}

// ErrUnknownState - will be returned wrapped when parsing string
// containing unrecognized value.
var ErrUnknownType = errors.New("unknown Type")

var mapTypeFromString = map[string]TaskType{
	"file": FileTask,
	"url":  URLTask,
}

// UnmarshalJSON implements the Unmarshaler interface of the json package for State.
func (t *TaskType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	result, ok := mapTypeFromString[strings.ToLower(v)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownType, v)
	}
	*t = result
	return nil
}

// MarshalJSON implements the Marshaler interface of the json package for State.
func (t TaskType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", t)), nil
}
