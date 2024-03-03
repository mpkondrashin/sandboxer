/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_intro.go

First installer page
*/
package main

import (
	"sandboxer/pkg/globals"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	reinstall = "Reinstall"
	uninstall = "Uninstall"
)

type PageReinstall struct {
	BasePage
	reinstallRadio *widget.RadioGroup
}

var _ Page = &PageReinstall{}

func (p *PageReinstall) Name() string {
	return "Reinstall"
}

func (p *PageReinstall) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	if p.reinstallRadio == nil {
		return pgInstallation
	}
	switch p.reinstallRadio.Selected {
	case reinstall:
		return pgInstallation
	case abort:
		return pgExit
	case uninstall:
		return pgUninstall
	}
	return 0 //&PageInstallation{		BasePage: BasePage{wiz: p.wiz},		//reinstallRadio: nil,	}
}

func (p *PageReinstall) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	titleLabel := widget.NewLabel(globals.AppName + " " + globals.Version + " version is already installed.")
	p.reinstallRadio = widget.NewRadioGroup([]string{abort, reinstall, uninstall}, p.radioChanged)
	p.reinstallRadio.SetSelected(abort)
	return container.NewVBox(
		titleLabel,
		p.reinstallRadio,
	)
}

func (p *PageReinstall) Run(win fyne.Window, installer *Installer) {

}

func (p *PageReinstall) AquireData(installer *Installer) error {
	switch p.reinstallRadio.Selected {
	case reinstall, uninstall:
		return nil
	case abort:
		return ErrAbort
	}
	return nil
}
func (p *PageReinstall) radioChanged(s string) {
	if p.reinstallRadio == nil {
		return
	}
	p.wiz.UpdatePagesList()
}
