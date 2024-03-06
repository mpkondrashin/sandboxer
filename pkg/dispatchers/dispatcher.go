/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

dispatcher.go

Base dispatcher functions
*/
package dispatchers

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"sandboxer/pkg/config"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/sandbox"
	"sandboxer/pkg/task"

	"github.com/google/uuid"
	"github.com/mpkondrashin/ddan"
	"github.com/mpkondrashin/vone"
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
	token := d.conf.VisionOne.Token
	if token == "" {
		return nil, errors.New("token is not set")
	}
	domain := d.conf.VisionOne.Domain
	if domain == "" {
		return nil, errors.New("domain is not set")
	}
	return sandbox.NewVOneSandbox(vone.NewVOne(domain, token)), nil
}

func (d *BaseDispatcher) AnalyzerSandbox() (sandbox.Sandbox, error) {
	u, err := url.Parse(d.conf.DDAn.URL)
	if err != nil {
		return nil, err
	}
	analyzer := ddan.NewClient(d.conf.DDAn.ProductName, d.conf.DDAn.Hostname)
	analyzer.SetAnalyzer(u, d.conf.DDAn.APIKey, d.conf.DDAn.IgnoreTLSErrors)
	//SetProtocolVersion(version string) ClientInterface
	analyzer.SetSource(d.conf.DDAn.SourceID, d.conf.DDAn.SourceName)
	if d.conf.DDAn.ClientUUID == "" {
		d.conf.DDAn.ClientUUID, err = d.GenerateUUID()
		if err != nil {
			return nil, err
		}
	}
	analyzer.SetUUID(d.conf.DDAn.ClientUUID)
	analyzer.SetProtocolVersion(d.conf.DDAn.ProtocolVersion)

	return sandbox.NewDDAnSandbox(analyzer), nil
}

func (d *BaseDispatcher) GenerateUUID() (string, error) {
	path, err := globals.AnalyzerClientUUIDFilePath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		logging.Errorf("Client UUID file read error: %v", err)
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		uuidData, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		err = os.WriteFile(path, []byte(uuidData.String()), 0600)
		if err != nil {
			return "", err
		}
		return uuidData.String(), nil
	}
	return string(data), nil
}
