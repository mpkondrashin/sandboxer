package dispatchers

import (
	"fmt"
	"sandboxer/pkg/task"
	"strings"
)

/*

[Submit] --- inbox ---> Prefilter ----> [Prefilter] --- Submit ---> [Submit] --- Result ---> [Result] --- meta ---> [Pull virus name]

*/

const (
	ChannelSize = 1000
)

const (
	ChPrefilter = iota
	ChUpload
	ChWait
	ChResult
	ChReport
	ChInvestigation
	ChCount
)

type IDChannel chan task.ID

type Channels struct {
	TaskChannel [ChCount]chan task.ID
}

func NewChannels() *Channels {
	c := &Channels{}
	for i := 0; i < ChCount; i++ {
		c.TaskChannel[i] = make(chan task.ID, ChannelSize)
	}
	return c
}

func (c *Channels) Close() {
	for i := 0; i < ChCount; i++ {
		close(c.TaskChannel[i])
	}
}

func (c *Channels) String() string {
	var sb strings.Builder
	for i := ChPrefilter; i < ChCount; i++ {
		fmt.Fprintf(&sb, "%d - ", len(c.TaskChannel[i]))
	}
	return sb.String()
}
