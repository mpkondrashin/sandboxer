package main

import (
	"examen/pkg/globals"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type PageLaunch struct {
	//tokenEntry *widget.Entry
	autostartCheck *widget.Check
}

var _ Page = &PageLaunch{}

func (p *PageLaunch) Name() string {
	return "Autostart"
}

func (p *PageLaunch) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	p.autostartCheck = widget.NewCheck("Add "+globals.AppName+" to autostart", nil)
	p.autostartCheck.SetChecked(installer.autostart)
	return p.autostartCheck
}

func (p *PageLaunch) Run(win fyne.Window, installer *Installer) {

}

func (p *PageLaunch) AquireData(installer *Installer) error {
	installer.autostart = p.autostartCheck.Checked
	//	installer.config.VisionOne.Token = p.tokenEntry.Text
	return nil
}
