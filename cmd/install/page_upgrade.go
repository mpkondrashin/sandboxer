/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_intro.go

First installer page
*/
package main

import (
	"fmt"
	"sandboxer/pkg/globals"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	//abort  = "Abort installation"
	upgrade = "Upgrade"
)

type PageUpgrade struct {
	BasePage
	reinstallRadio *widget.RadioGroup
	content        fyne.CanvasObject
}

var _ Page = &PageUpgrade{}

func (p *PageUpgrade) Name() string {
	return "Upgrade"
}

func (p *PageUpgrade) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	if p.reinstallRadio == nil {
		return pgInstallation
	}
	switch p.reinstallRadio.Selected {
	case upgrade:
		return pgInstallation
	case abort:
		return pgExit
	case uninstall:
		return pgUninstall //&PageUninstall{BasePage{wiz: p.wiz}}
	}
	return pgFinish
	//p.w.UpdatePagesList(p)

	//return &PageFinish{BasePage{wiz: p.wiz}}
}

func (p *PageUpgrade) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	if p.content != nil {
		return p.content
	}
	titleLabel := widget.NewRichTextWithText(fmt.Sprintf(globals.AppName+" %s version is already installed on this system. Upgrade to %s?",
		installer.config.Version, globals.Version))
	titleLabel.Wrapping = fyne.TextWrapWord
	p.reinstallRadio = widget.NewRadioGroup([]string{upgrade, abort, uninstall}, p.radioChanged)
	p.reinstallRadio.SetSelected(upgrade)
	p.content = container.NewVBox(
		titleLabel,
		p.reinstallRadio,
	)
	return p.content

}

func (p *PageUpgrade) Run(win fyne.Window, installer *Installer) {

}

func (p *PageUpgrade) AquireData(installer *Installer) error {
	switch p.reinstallRadio.Selected {
	case reinstall, uninstall:
		return nil
	case abort:
		return ErrAbort
	}
	return nil
}
func (p *PageUpgrade) radioChanged(s string) {
	if p.reinstallRadio == nil {
		return
	}
	p.wiz.UpdatePagesList()
	// update page
}
