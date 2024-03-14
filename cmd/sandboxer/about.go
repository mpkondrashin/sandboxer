/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

about.go

About window
*/
package main

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/globals"
)

type AboutWindow struct {
}

func NewAboutWindow() *AboutWindow {
	return &AboutWindow{}
}

func (w *AboutWindow) Content(modal *ModalWindow) fyne.CanvasObject {
	name := widget.NewLabel(globals.AppName)
	version := widget.NewLabel(fmt.Sprintf("Version %s Build %s", globals.Version, globals.Build))
	repoURL, _ := url.Parse("https://github.com/mpkondrashin/" + globals.Name)
	repoLink := widget.NewHyperlink("Repository on GitHub", repoURL)
	vbox := container.NewVBox(
		container.NewCenter(name),
		container.NewCenter(version),
		container.NewCenter(repoLink),
	)
	return container.NewPadded(vbox)
}

func (s *AboutWindow) Name() string {
	return "About"
}

func (s *AboutWindow) Icon() fyne.Resource {
	return theme.InfoIcon()
}

func (s *AboutWindow) Show() {}
func (s *AboutWindow) Hide() {}
