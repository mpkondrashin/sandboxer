package main

import (
	"examen/pkg/globals"
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	IntoText = "Examen provides ability to check files using Vision One sandbox. " +
		"You will have to provide API Key to use this service."

	License = `MIT License

Copyright (c) 2024 Michael Kondrashin (mkondrashin@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`
)

type PageIntro struct {
}

var _ Page = &PageIntro{}

func (p *PageIntro) Name() string {
	return "Intro"
}

func (p *PageIntro) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	titleLabel := widget.NewLabelWithStyle(globals.AppName,
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	version := fmt.Sprintf("Version %s", MajorMinorRevision[1:])
	versionLabel := widget.NewLabelWithStyle(version,
		fyne.TextAlignCenter, fyne.TextStyle{})

	report := widget.NewRichTextFromMarkdown(IntoText)
	report.Wrapping = fyne.TextWrapWord
	repoURL, _ := url.Parse("https://github.com/mpkondrashin/examen")
	repoLink := widget.NewHyperlink("Examen repository on GitHub", repoURL)
	coneURL, _ := url.Parse("https://portal.xdr.trendmicro.com")
	coneLink := widget.NewHyperlink("Open Vision One Console", coneURL)

	licensePopUp := func() {
		licenseLabel := widget.NewLabel(License)
		sc := container.NewScroll(licenseLabel)
		popup := dialog.NewCustom("Show License Information", "Close", sc, win)
		popup.Resize(fyne.NewSize(800, 600))
		popup.Show()
	}
	licenseButton := widget.NewButton("License Information...", licensePopUp)
	return container.NewVBox(
		titleLabel,
		versionLabel,
		report,
		container.NewHBox(coneLink),
		container.NewHBox(repoLink, licenseButton),
	)
}

func (p *PageIntro) AquireData(installer *Installer) error {
	return nil
}
