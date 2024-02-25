/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

sandbox_type.go

Type of supported sandboxes
*/
package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SandboxType int

const (
	SandboxVisionOne SandboxType = iota
	SandboxAnalyzer
)

// String - return string representation for State value
func (v SandboxType) String() string {
	s, ok := map[SandboxType]string{
		SandboxVisionOne: "VisionOne",
		SandboxAnalyzer:  "Analyzer",
	}[v]
	if ok {
		return s
	}
	return "SandboxType(" + strconv.FormatInt(int64(v), 10) + ")"
}

// ErrUnknownState - will be returned wrapped when parsing string
// containing unrecognized value.
var ErrUnknownState = errors.New("unknown State")

var mapSandboxTypeFromString = map[string]SandboxType{
	"visionone": SandboxVisionOne,
	"analyzer":  SandboxAnalyzer,
}

// UnmarshalJSON implements the Unmarshaler interface of the json package for State.
func (s *SandboxType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	result, ok := mapSandboxTypeFromString[strings.ToLower(v)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownState, v)
	}
	*s = result
	return nil
}

// MarshalJSON implements the Marshaler interface of the json package for State.
func (s SandboxType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", s)), nil
}

// UnmarshalYAML implements the Unmarshaler interface of the yaml.v3 package for State.
func (s *SandboxType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	result, ok := mapSandboxTypeFromString[strings.ToLower(v)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownState, v)
	}
	*s = result
	return nil
}

// MarshalXML implements the Marshaler interface of the xml package for State.
func (s SandboxType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(s.String(), start)
}

// UnmarshalXML implements the Unmarshaler interface of the xml package for State.
func (s *SandboxType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	result, ok := mapSandboxTypeFromString[strings.ToLower(v)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownState, v)
	}
	*s = result
	return nil
}
