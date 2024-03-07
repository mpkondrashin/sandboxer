/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

dispatcher.go

Base dispatcher functions
*/
package dispatchers

import (
	"fmt"
	"sandboxer/pkg/config"
	"sandboxer/pkg/sandbox"
	"sandboxer/pkg/task"
)

type Dispatcher interface {
	InboundChannel() task.Channel
	ProcessTask(tsk *task.Task) error
}

type BaseDispatcher struct {
	conf     *config.Configuration
	channels *task.Channels
	list     *task.TaskList
}

func NewBaseDispatcher(conf *config.Configuration, channels *task.Channels, list *task.TaskList) BaseDispatcher {
	return BaseDispatcher{conf, channels, list}
}

func (d *BaseDispatcher) Channel(ch task.Channel) task.IDChannel {
	return d.channels.TaskChannel[ch]
}

func (d *BaseDispatcher) Sandbox() (sandbox.Sandbox, error) {
	switch d.conf.SandboxType {
	case config.SandboxVisionOne:
		return d.VisionOneSandbox()
	case config.SandboxAnalyzer:
		return d.AnalyzerSandbox()
	}
	return nil, fmt.Errorf("uknown Sandbox Type: %d", d.conf.SandboxType)
}

func (d *BaseDispatcher) VisionOneSandbox() (sandbox.Sandbox, error) {
	vOne, err := d.conf.VisionOne.VisionOneSandbox()
	if err != nil {
		return nil, err
	}
	return sandbox.NewVOneSandbox(vOne), nil
}

func (d *BaseDispatcher) AnalyzerSandbox() (sandbox.Sandbox, error) {
	analyzer, err := d.conf.DDAn.AnalyzerWithUUID()
	if err != nil {
		return nil, err
	}
	return sandbox.NewDDAnSandbox(analyzer), nil
}
