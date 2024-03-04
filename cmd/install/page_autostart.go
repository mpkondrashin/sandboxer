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
	"sandboxer/pkg/xplatform"
)

type PageAutostart struct {
	BasePage
	autostartCheck *widget.Check
}

var _ Page = &PageAutostart{}

func (p *PageAutostart) Name() string {
	return "Autostart"
}

func (p *PageAutostart) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	if xplatform.IsWindows() {
		return pgFolder
	}
	return pgInstallation
}

func (p *PageAutostart) Content() fyne.CanvasObject {
	p.autostartCheck = widget.NewCheck("Add "+globals.AppName+" to autostart", nil)
	p.autostartCheck.SetChecked(p.wiz.installer.autostart)
	return p.autostartCheck
}

//func (p *PageAutostart) Run() {}

func (p *PageAutostart) AquireData(installer *Installer) error {
	installer.autostart = p.autostartCheck.Checked
	return nil
}
