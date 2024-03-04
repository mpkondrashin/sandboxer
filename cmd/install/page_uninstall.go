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

type PageUninstall struct {
	BasePage
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
}

var _ Page = &PageUninstall{}

func (p *PageUninstall) Name() string {
	return "Uninstall"
}

func (p *PageUninstall) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgExit
}

func (p *PageUninstall) Content() fyne.CanvasObject {
	p.progressBar = widget.NewProgressBar()
	p.statusLabel = widget.NewLabel("")
	return container.NewVBox(
		p.progressBar,
		p.statusLabel,
		//copyButton,
	)
}

func (p *PageUninstall) Run() {
	total := float64(len(p.wiz.installer.UninstallStages()) - 1)
	index := 0
	stageName := ""
	err := p.wiz.installer.Uninstall(func(name string) error {
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
		dialog.ShowError(err, p.wiz.win)
	}
}

func (p *PageUninstall) AquireData(installer *Installer) error {
	return nil
}
