/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_token.go

Provide Vision One token and domain
*/
package main

import (
	"fyne.io/fyne/v2"

	"sandboxer/pkg/config"
	"sandboxer/pkg/settings"
)

type PageProxy struct {
	BasePage
	proxySettings *settings.Proxy
}

var _ Page = &PageProxy{}

func (p *PageProxy) Name() string {
	return "Proxy"
}

func (p *PageProxy) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	switch p.wiz.installer.config.SandboxType {
	case config.SandboxVisionOne:
		return pgVOSettings
	case config.SandboxAnalyzer:
		return pgDDSettings
	}
	return pgVOSettings
}

func (p *PageProxy) Content() fyne.CanvasObject {
	p.proxySettings = settings.NewProxy(p.wiz.installer.config.Proxy)
	return p.proxySettings.Widget()
}

func (p *PageProxy) Run() {
	// No need to load, config is loaded when application started
	//	err := installer.LoadConfig()
	//	if err != nil {
	//		logging.Errorf("LoadConfig: %v", err)
	//		dialog.ShowError(err, win)
	//	}
}

func (p *PageProxy) AquireData(installer *Installer) error {
	err := p.proxySettings.Aquire()
	if err != nil {
		return err
	}
	*installer.config.Proxy = *p.proxySettings.Conf
	return nil
}
