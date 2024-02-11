package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/logging"
)

type PageOptions struct {
	tokenEntry *widget.Entry
}

var _ Page = &PageOptions{}

func (p *PageOptions) Name() string {
	return "Options"
}

func (p *PageOptions) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
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

func (p *PageOptions) Run(win fyne.Window, installer *Installer) {
	err := installer.LoadConfig()
	if err != nil {
		logging.Errorf("LoadConfig: %v", err)
		dialog.ShowError(err, win)
	}
	p.tokenEntry.SetText(installer.config.VisionOne.Token)

}

func (p *PageOptions) AquireData(installer *Installer) error {
	if p.tokenEntry.Text == "" {
		return fmt.Errorf("token field is empty")
	}
	installer.config.VisionOne.Token = p.tokenEntry.Text
	return nil
}
