package main

import (
	"examen/pkg/logging"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageOptions struct {
	tokenEntry *widget.Entry
}

var _ Page = &PageOptions{}

func (p *PageOptions) Name() string {
	return "Options"
}

func (p *PageOptions) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	err := installer.LoadConfig()
	if err != nil {
		logging.Errorf("LoadConfig: %v", err)
		dialog.ShowError(err, win)
	}
	labelTop := widget.NewLabel("Please open Vision One console to get all nessesary parameters")
	p.tokenEntry = widget.NewMultiLineEntry()
	p.tokenEntry.Text = installer.config.Token
	p.tokenEntry.Wrapping = fyne.TextWrapBreak
	tokenFormItem := widget.NewFormItem("Token:", p.tokenEntry)
	tokenFormItem.HintText = "Go to XXXXXXX"
	optionsForm := widget.NewForm(
		tokenFormItem,
	)
	return container.NewVBox(labelTop, optionsForm)
}

func (p *PageOptions) AquireData(installer *Installer) error {
	if p.tokenEntry.Text == "" {
		return fmt.Errorf("token field is empty")
	}
	installer.config.Token = p.tokenEntry.Text
	return nil
}
