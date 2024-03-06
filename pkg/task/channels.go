/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

channels.go

Channels to dispatchers to communicate
*/
package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	ChannelSize = 1000
)

type Channel int

const (
	ChPrefilter Channel = iota
	ChSubmit
	//	ChWait
	ChResult
	ChReport
	ChInvestigation
	ChDone
)

var ChannelString = [...]string{
	"Prefilter",
	"Submission",
	//	"Wait",
	"Wait For Result",
	"Get Report",
	"Get Investigation",
	"Done",
}

// String - return string representation for State value
func (c Channel) String() string {
	if c < 0 || c > ChDone {
		return "Channel(" + strconv.FormatInt(int64(c), 10) + ")"
	}
	return ChannelString[c]
}

// ErrUnknownState - will be returned wrapped when parsing string
// containing unrecognized value.
var ErrUnknownChannel = errors.New("unknown Channel")

var mapChannelFromString = map[string]Channel{
	ChannelString[ChPrefilter]: ChPrefilter,
	ChannelString[ChSubmit]:    ChSubmit,
	//	ChannelString[ChWait]:          ChWait,
	ChannelString[ChResult]:        ChResult,
	ChannelString[ChReport]:        ChReport,
	ChannelString[ChInvestigation]: ChInvestigation,
	ChannelString[ChDone]:          ChDone,
	ChannelString[ChPrefilter]:     ChPrefilter,
}

// UnmarshalJSON implements the Unmarshaler interface of the json package for State.
func (c *Channel) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	result, ok := mapChannelFromString[v]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownChannel, v)
	}
	*c = result
	return nil
}

// MarshalJSON implements the Marshaler interface of the json package for State.
func (c Channel) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", c)), nil
}

type IDChannel chan ID

type Channels struct {
	TaskChannel [ChDone]chan ID
}

func NewChannels() *Channels {
	c := &Channels{}
	for i := ChPrefilter; i < ChDone; i++ {
		c.TaskChannel[i] = make(chan ID, ChannelSize)
	}
	return c
}

func (c *Channels) Close() {
	for i := ChPrefilter; i < ChDone; i++ {
		close(c.TaskChannel[i])
	}
}

func (c *Channels) String() string {
	var sb strings.Builder
	for i := ChPrefilter; i < ChDone; i++ {
		fmt.Fprintf(&sb, "%d - ", len(c.TaskChannel[i]))
	}
	return sb.String()
}
