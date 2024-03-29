/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_intro.go

First installer page
*/
package main

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/config"
	"sandboxer/pkg/globals"
)

const (
	IntoText = globals.AppName + " provides ability to check files using Vision One sandbox service or Deep Discovery Analyzer appliance."

	NoteText = "Please close all MMC windows before continuing."
)

type PageIntro struct {
	BasePage
	sandboxRadio *widget.RadioGroup
	proxyCheck   *widget.Check
}

var _ Page = &PageIntro{}

func (p *PageIntro) Name() string {
	return "Intro"
}

func (p *PageIntro) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	if p.sandboxRadio == nil {
		return pgVOSettings
	}
	if p.proxyCheck != nil && p.proxyCheck.Checked {
		return pgProxy
	}
	switch p.sandboxRadio.Selected {
	case SandboxVisionOne:
		return pgVOSettings
	case SandboxDDAn:
		return pgDDSettings
	default:
		return pgExit
	}
	//return pgVOSettings
}

const (
	SandboxVisionOne = "Vision One Sandbox Service"
	SandboxDDAn      = "Deep Discovery Analyzer"
)

func (p *PageIntro) Content() fyne.CanvasObject {
	titleLabel := widget.NewLabelWithStyle(globals.AppName,
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	version := fmt.Sprintf("Version %s build %s", globals.Version, globals.Build)
	versionLabel := widget.NewLabelWithStyle(version,
		fyne.TextAlignCenter, fyne.TextStyle{})

	report := widget.NewRichTextFromMarkdown(IntoText)
	report.Wrapping = fyne.TextWrapWord

	chooseLabel := widget.NewLabel("Choose your sandbox:")
	p.sandboxRadio = widget.NewRadioGroup(
		[]string{SandboxVisionOne, SandboxDDAn},
		nil,
	)
	p.sandboxRadio.Required = true
	switch p.wiz.installer.config.SandboxType {
	case config.SandboxVisionOne:
		p.sandboxRadio.SetSelected(SandboxVisionOne)
	case config.SandboxAnalyzer:
		p.sandboxRadio.SetSelected(SandboxDDAn)
	}

	p.proxyCheck = widget.NewCheck("Use proxy", p.proxyChanged)
	p.proxyCheck.Checked = p.wiz.installer.config.Proxy.Active

	noteMarkdown := widget.NewRichTextFromMarkdown(NoteText)
	noteMarkdown.Wrapping = fyne.TextWrapWord

	repoURL, _ := url.Parse("https://github.com/mpkondrashin/" + globals.Name)
	repoLink := widget.NewHyperlink(globals.AppName+" repository on GitHub", repoURL)

	licensePopUp := func() {
		licenseLabel := widget.NewLabel(LicenseText())
		sc := container.NewScroll(licenseLabel)
		popup := dialog.NewCustom("Show License Information", "Close", sc, p.wiz.win)
		popup.Resize(fyne.NewSize(800, 600))
		popup.Show()
	}
	licenseButton := widget.NewButton("License Information...", licensePopUp)
	return container.NewVBox(
		titleLabel,
		versionLabel,
		report,
		chooseLabel,
		p.sandboxRadio,
		p.proxyCheck,
		noteMarkdown,
		container.NewHBox(repoLink, licenseButton),
	)
}

func (p *PageIntro) Run() {
	fmt.Println("Run" + p.Name())
	fmt.Println("Type ", p.wiz.installer.config.SandboxType.String())
	//p.sandboxRadio.SetSelected(p.wiz.installer.config.SandboxType.String())
}

func (p *PageIntro) proxyChanged(checked bool) {
	if p.proxyCheck == nil {
		return
	}
	p.wiz.UpdatePagesList()
}

func (p *PageIntro) AquireData(installer *Installer) error {
	switch p.sandboxRadio.Selected {
	case SandboxVisionOne:
		installer.config.SandboxType = config.SandboxVisionOne
	case SandboxDDAn:
		installer.config.SandboxType = config.SandboxAnalyzer
	default:
		return errors.New("Choose the Sandbox to be used")
	}
	p.wiz.installer.config.Proxy.Active = p.proxyCheck.Checked
	return nil
}

func LicenseText() string {
	filePath := "embed/LICENSE"
	licFile, err := embedFS.Open(filePath)
	if err != nil {
		return "reading error"
	}
	defer func() {
		licFile.Close()
	}()
	licBytes, err := io.ReadAll(licFile)
	if err != nil {
		return "reading error"
	}
	return string(licBytes)

}
