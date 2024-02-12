/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_autostart.go

Checkbox for autostart
*/
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/globals"
)

type PageAutostart struct {
	autostartCheck *widget.Check
}

var _ Page = &PageAutostart{}

func (p *PageAutostart) Name() string {
	return "Autostart"
}

func (p *PageAutostart) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	p.autostartCheck = widget.NewCheck("Add "+globals.AppName+" to autostart", nil)
	p.autostartCheck.SetChecked(installer.autostart)
	return p.autostartCheck
}

func (p *PageAutostart) Run(win fyne.Window, installer *Installer) {}

func (p *PageAutostart) AquireData(installer *Installer) error {
	installer.autostart = p.autostartCheck.Checked
	return nil
}
