package main

import (
	"examen/pkg/globals"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type PageFinish struct {
}

var _ Page = &PageAutostart{}

func (p *PageFinish) Name() string {
	return "Finish"
}

var finalText = "Examen service sucessfully installed.\n\n" +
	"Right click on any file and pick Send To -> " + globals.AppName + "."

func (p *PageFinish) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	report := widget.NewRichTextFromMarkdown(finalText)
	report.Wrapping = fyne.TextWrapWord
	return report
}

func (p *PageFinish) Run(win fyne.Window, installer *Installer) {}

func (p *PageFinish) AquireData(installer *Installer) error {
	return nil
}
