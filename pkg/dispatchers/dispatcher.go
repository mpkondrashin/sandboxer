package dispatchers

import (
	"errors"
	"sandboxer/pkg/config"
	"sandboxer/pkg/task"

	"github.com/mpkondrashin/vone"
)

type Dispatcher interface {
	InboundChannel() int
	ProcessTask(tsk *task.Task) error
}

type BaseDispatcher struct {
	conf     *config.Configuration
	channels *Channels
	list     *task.TaskList
}

func NewBaseDispatcher(conf *config.Configuration, channels *Channels, list *task.TaskList) BaseDispatcher {
	return BaseDispatcher{conf, channels, list}
}

func (d *BaseDispatcher) Channel(ch int) IDChannel {
	return d.channels.TaskChannel[ch]
}

func (d *BaseDispatcher) vOne() (*vone.VOne, error) {
	token := d.conf.VisionOne.Token
	if token == "" {
		return nil, errors.New("token is not set")
	}
	domain := d.conf.VisionOne.Domain
	if domain == "" {
		return nil, errors.New("domain is not set")
	}
	return vone.NewVOne(domain, token), nil
}
