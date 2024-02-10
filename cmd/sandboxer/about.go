package main

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/globals"
)

type AboutWindow struct {
	ModalWindow
}

func NewAboutWindow(modalWindow ModalWindow) *AboutWindow {
	s := &AboutWindow{
		ModalWindow: modalWindow,
	}
	name := widget.NewLabel(globals.AppName)
	version := widget.NewLabel(fmt.Sprintf("Version %s Build %s", globals.Version, globals.Build))
	repoURL, _ := url.Parse("https://github.com/mpkondrashin/" + globals.Name)
	repoLink := widget.NewHyperlink("Repository on GitHub", repoURL)

	vbox := container.NewVBox(
		container.NewCenter(name),
		container.NewCenter(version),
		container.NewCenter(repoLink),
	)
	s.win.SetContent(vbox)
	return s
}

func (s *AboutWindow) Show() {
	s.win.Show()
}
