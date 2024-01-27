package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/globals"
)

type PageFinish struct {
}

var _ Page = &PageAutostart{}

func (p *PageFinish) Name() string {
	return "Finish"
}

func (p *PageFinish) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	l1 := widget.NewLabel(globals.AppName + " service sucessfully installed.")
	l2 := widget.NewLabel("Right click on any file and pick Send To -> " + globals.AppName + ".")
	return container.NewVBox(l1, l2)
}

func (p *PageFinish) Run(win fyne.Window, installer *Installer) {}

func (p *PageFinish) AquireData(installer *Installer) error {
	return nil
}
