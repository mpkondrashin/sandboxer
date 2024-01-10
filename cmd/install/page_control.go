package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageInstallation struct {
}

var _ Page = &PageOptions{}

func (p *PageInstallation) Name() string {
	return "Control"
}

func (p *PageInstallation) GetStatus(model *Model) {
}

func (p *PageInstallation) Content(win fyne.Window, model *Model) fyne.CanvasObject {
	progressBar := widget.NewProgressBar()
	var copyButton *widget.Button
	copyButton = widget.NewButton("Copy Files", func() {
		copyButton.Disable()
		for i := 0; i < 10; i++ {
			progressBar.SetValue(float64(i) / 10.)
			time.Sleep(1 * time.Second)
		}
	})
	return container.NewVBox(
		progressBar,
		copyButton,
	)
}

func (p *PageInstallation) AquireData(model *Model) error {
	return nil
}
