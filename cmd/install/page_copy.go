/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_autostart.go

Copy files
*/
package main

import (
	"fmt"
	"sandboxer/pkg/logging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageInstallation struct {
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
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
	p.progressBar = widget.NewProgressBar()
	p.statusLabel = widget.NewLabel("")
	//	var copyButton *widget.Button
	//copyButton = widget.NewButton("Copy Files",
	//func() {
	//	copyButton.Disable()

	//})
	return container.NewVBox(
		p.progressBar,
		p.statusLabel,
		//copyButton,
	)
}

func (p *PageInstallation) Run(win fyne.Window, installer *Installer) {
	total := float64(len(installer.Stages()) - 1)
	index := 0
	stageName := ""
	err := installer.Install(func(name string) error {
		stageName = name
		p.progressBar.SetValue(float64(index) / total)
		p.statusLabel.SetText(name)
		index++
		return nil
	})
	p.statusLabel.SetText("Done")
	if err != nil {
		p.statusLabel.SetText(stageName + " Failed")
		err = fmt.Errorf("%s: %w", stageName, err)
		logging.LogError(err)
		dialog.ShowError(err, win)
	}
}

func (p *PageInstallation) AquireData(installer *Installer) error {
	return nil
}
