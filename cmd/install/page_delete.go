/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_intro.go

First installer page
*/
package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	abort  = "Abort installation"
	delete = "Delete existing config and install"
)

type PageDelete struct {
	BasePage
	ErrorMessage string
	deleteRadio  *widget.RadioGroup
}

var _ Page = &PageDelete{}

func (p *PageDelete) Name() string {
	return "Error"
}

func (p *PageDelete) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	if p.deleteRadio == nil {
		return pgExit
	}
	switch p.deleteRadio.Selected {
	case delete:
		return pgIntro
	case abort:
		return pgExit
	}
	return pgIntro
}

func (p *PageDelete) Content() fyne.CanvasObject {
	titleLabel := widget.NewLabel("Configuration Error")

	errorText := widget.NewRichTextFromMarkdown("Error Message: " + p.ErrorMessage)
	errorText.Wrapping = fyne.TextWrapWord

	p.deleteRadio = widget.NewRadioGroup([]string{abort, delete}, p.radioChanged)
	p.deleteRadio.SetSelected(abort)
	return container.NewVBox(
		titleLabel,
		errorText,
		p.deleteRadio,
	)
}

func (p *PageDelete) Run() {

}

func (p *PageDelete) AquireData(installer *Installer) error {
	switch p.deleteRadio.Selected {
	case delete:
		return os.Remove(installer.config.GetFilePath())
	case abort:
		return ErrAbort
	}
	return nil
}

func (p *PageDelete) radioChanged(s string) {
	if p.deleteRadio == nil {
		return
	}
	p.wiz.UpdatePagesList()
	// update page
}
