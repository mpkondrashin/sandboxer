/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_token.go

Provide Vision One token and domain
*/
package main

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/settings"
)

type PageVOToken struct {
	BasePage
	voneSettings *settings.VisionOne
}

var _ Page = &PageVOToken{}

func (p *PageVOToken) Name() string {
	return "Settings"
}

func (p *PageVOToken) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgAutostart
}

func (p *PageVOToken) Content() fyne.CanvasObject {
	p.voneSettings = settings.NewVisionOne(p.wiz.installer.config.VisionOne)

	//labelTop := widget.NewLabel("Open Vision One console to get all nessesary parameters")

	voneURL, _ := url.Parse("https://portal.xdr.trendmicro.com")
	voneLink := widget.NewHyperlink("Open Vision One console to get all nessesary parameters", voneURL)

	// https://docs.trendmicro.com/en-US/documentation/article/trend-vision-one-configuring-user-rol
	// https://docs.trendmicro.com/en-us/documentation/article/trend-vision-one-api-keys
	return container.NewVBox(voneLink, p.voneSettings.Widget())
}

func (p *PageVOToken) Run() {
	// No need to load, config is loaded when application started
	//	err := installer.LoadConfig()
	//	if err != nil {
	//		logging.Errorf("LoadConfig: %v", err)
	//		dialog.ShowError(err, win)
	//	}
	fmt.Println("Run" + p.Name())
}

func (p *PageVOToken) AquireData(installer *Installer) error {
	err := p.voneSettings.Aquire()
	if err != nil {
		return err
	}
	installer.config.VisionOne = p.voneSettings.Conf
	return nil
}
