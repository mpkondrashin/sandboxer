/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

quota.go

Quota window
*/
package main

import (
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/mod/semver"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/update"
)

type UpdateWindow struct {
	ModalWindow
	//conf    *config.Configuration
	//version        binding.String
	version        string
	versionLabel   *widget.Label
	progressBar    *widget.ProgressBar
	downloadButton *widget.Button
}

func NewUpdateWindow(modalWindow ModalWindow /*, conf *config.Configuration*/) *UpdateWindow {
	s := &UpdateWindow{
		ModalWindow: modalWindow,
		//	conf:        conf,
		//version:     binding.NewString(),
		versionLabel: widget.NewLabel("                                 "),
		progressBar:  widget.NewProgressBar(),
	}
	s.downloadButton = widget.NewButton("Download", s.Download)
	//s.progressBar
	s.downloadButton.Disable()
	s.Reset()
	//versionCountItem := widget.NewFormItem("Latest Available Version:", widget.NewLabelWithData(s.version))
	//form := widget.NewForm(
	//versionCountItem,
	//)
	s.win.SetContent(container.NewVBox(s.versionLabel, s.progressBar, s.downloadButton))
	return s
}

func (s *UpdateWindow) Download() {
	s.downloadButton.Disable()
	fileName := "Setup.zip"
	if runtime.GOOS == "darwin" {
		fileName = "Setup_macos.zip"
	}

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
	//s.version.Set("?         ")
	s.versionLabel.SetText("Checking...")
	s.version = ""
}

func (s *UpdateWindow) Update() {
	s.Reset()
	//exemptionCount := fmt.Sprintf("%d (Number of samples submitted but marked as \"not analyzed\". This number does not count toward the daily reserve)", result.SubmissionExemptionCount)
	var err error
	s.version, err = update.LatestVersion(globals.Name)
	if err != nil {
		dialog.ShowError(err, s.win)
		return
	}
	//s.version.Set(version)
	switch semver.Compare(s.version, globals.Version) {
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
