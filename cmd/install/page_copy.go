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
	statusLabel := widget.NewLabel("")
	var copyButton *widget.Button
	copyButton = widget.NewButton("Copy Files", func() {
		copyButton.Disable()
		total := float64(len(installer.Stages()) - 1)
		index := 0
		err := installer.Install(func(name string) error {
			progressBar.SetValue(float64(index) / total)
			statusLabel.SetText(name)
			index++
			return nil
		})
		logging.LogError(err)
		if err != nil {
			dialog.ShowError(err, win)
		}
	})
	return container.NewVBox(
		progressBar,
		statusLabel,
		copyButton,
	)
}

func (p *PageInstallation) AquireData(installer *Installer) error {
	return nil
}
