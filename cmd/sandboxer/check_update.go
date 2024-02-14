/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

quota.go

Quota window
*/
package main

import (
	"fmt"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/mod/semver"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/update"
)

type UpdateWindow struct {
	ModalWindow
	version        string
	versionLabel   *widget.Label
	progressBar    *widget.ProgressBar
	downloadButton *widget.Button
}

func NewUpdateWindow(modalWindow ModalWindow) *UpdateWindow {
	s := &UpdateWindow{
		ModalWindow:  modalWindow,
		versionLabel: widget.NewLabel("                                 "),
		progressBar:  widget.NewProgressBar(),
	}
	s.downloadButton = widget.NewButton("Download", s.Download)
	s.downloadButton.Disable()
	s.Reset()
	s.win.SetContent(container.NewPadded(container.NewVBox(s.versionLabel, s.progressBar, s.downloadButton)))
	return s
}

func (s *UpdateWindow) Download() {
	s.downloadButton.Disable()
	fileName := fmt.Sprintf("setup_%s_%s.zip", runtime.GOOS, runtime.GOARCH)
	update.DownloadRelease(s.version, fileName, globals.DownloadsFolder(), func(p float32) error {
		s.progressBar.SetValue(float64(p))
		return nil
	})
	err := RunOpen(filepath.Join(globals.DownloadsFolder(), fileName))
	if err != nil {
		dialog.ShowError(err, s.win)
	}
}

func (s *UpdateWindow) Reset() {
	s.versionLabel.SetText("Checking...")
	s.version = ""
}

func (s *UpdateWindow) Update() {
	s.Reset()
	var err error
	s.version, err = update.LatestVersion(globals.Name)
	if err != nil {
		dialog.ShowError(err, s.win)
		return
	}
	cmp := semver.Compare(s.version, globals.Version)
	logging.Debugf("Compare %s vs %s: %d", s.version, globals.Version, cmp)
	switch cmp {
	case -1:
	case 0:
		s.versionLabel.SetText("You have the newest version")
	case 1:
		s.versionLabel.SetText("New version available: " + s.version)
		s.downloadButton.Enable()
	}
}

func (s *UpdateWindow) Show() {
	s.win.Show()
	go s.Update()
}
