/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_token.go

Provide Vision One token
*/
package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageVOToken struct {
	BasePage
	tokenEntry *widget.Entry
}

var _ Page = &PageVOToken{}

func (p *PageVOToken) Name() string {
	return "Token"
}

func (p *PageVOToken) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgVODomain
}

func (p *PageVOToken) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	labelTop := widget.NewLabel("Please open Vision One console to get all nessesary parameters")
	p.tokenEntry = widget.NewMultiLineEntry()
	p.tokenEntry.Wrapping = fyne.TextWrapBreak
	tokenFormItem := widget.NewFormItem("Token:", p.tokenEntry)
	tokenFormItem.HintText = "Go to Administrator -> API Keys"
	// https://docs.trendmicro.com/en-US/documentation/article/trend-vision-one-configuring-user-rol
	// https://docs.trendmicro.com/en-us/documentation/article/trend-vision-one-api-keys
	optionsForm := widget.NewForm(
		tokenFormItem,
	)
	return container.NewVBox(labelTop, optionsForm)
}

func (p *PageVOToken) Run(win fyne.Window, installer *Installer) {
	// No need to load, config is loaded when application started
	//	err := installer.LoadConfig()
	//	if err != nil {
	//		logging.Errorf("LoadConfig: %v", err)
	//		dialog.ShowError(err, win)
	//	}
	p.tokenEntry.SetText(installer.config.VisionOne.Token)
}

func (p *PageVOToken) AquireData(installer *Installer) error {
	if p.tokenEntry.Text == "" {
		return fmt.Errorf("token field is empty")
	}
	installer.config.VisionOne.Token = strings.TrimSpace(p.tokenEntry.Text)
	return nil
}
