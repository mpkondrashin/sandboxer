package main

import (
	"examen/pkg/logging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageInstallation struct {
}

var _ Page = &PageInstallation{}

func (p *PageInstallation) Name() string {
	return "Copy Files"
}

/*
func (p *PageInstallation) GetStatus(installer *Installer) {
}
*/
func (p *PageInstallation) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	progressBar := widget.NewProgressBar()
	var copyButton *widget.Button
	stages := installer.Stages()

	copyButton = widget.NewButton("Copy Files", func() {
		copyButton.Disable()
		for i, stage := range stages {
			progressBar.SetValue(float64(i) / float64(len(stages)))
			err := stage()
			logging.LogError(err)
			if err != nil {
				dialog.ShowError(err, win)
				break
			}
		}
	})
	return container.NewVBox(
		progressBar,
		copyButton,
	)
}

func (p *PageInstallation) AquireData(installer *Installer) error {
	return nil
}
