/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_ddan.go

Provide Analyzer parameters
*/
package main

import (
	"sandboxer/pkg/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageDDAnSettings struct {
	BasePage
	ddanSettings *settings.DDAn
}

var _ Page = &PageDDAnSettings{}

func (p *PageDDAnSettings) Name() string {
	return "Settings"
}

func (p *PageDDAnSettings) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgAutostart
}

func (p *PageDDAnSettings) Content() fyne.CanvasObject {
	p.ddanSettings = settings.NewDDAnSettings(&p.wiz.installer.config.DDAn)

	labelTop := widget.NewLabel("Please provide your Deep Discovery Analyzer parameters")

	// https://docs.trendmicro.com/en-US/documentation/article/trend-vision-one-configuring-user-rol
	// https://docs.trendmicro.com/en-us/documentation/article/trend-vision-one-api-keys
	return container.NewVBox(labelTop, p.ddanSettings.Widget())
}

func (p *PageDDAnSettings) Run() {
	// No need to load, config is loaded when application started
	//	err := installer.LoadConfig()
	//	if err != nil {
	//		logging.Errorf("LoadConfig: %v", err)
	//		dialog.ShowError(err, win)
	//	}
	//p.tokenEntry.SetText(p.wiz.installer.config.VisionOne.Token)
}

func (p *PageDDAnSettings) AquireData(installer *Installer) error {
	return p.ddanSettings.Aquire()
}
