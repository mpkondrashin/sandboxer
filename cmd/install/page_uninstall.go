/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_uninstall.go

Run uninstallation
*/
package main

import (
	"fmt"
	"sandboxer/pkg/logging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageUninstall struct {
	BasePage
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
	errorsLabel *widget.Label
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
	p.errorsLabel = widget.NewLabel("")
	scroll := container.NewScroll(p.errorsLabel)
	scroll.SetMinSize(fyne.NewSize(0, 10*10))
	return container.NewVBox(
		p.progressBar,
		p.statusLabel,
		scroll,
	)
}

func (p *PageUninstall) Run() {
	stages := p.wiz.installer.UninstallStages()
	total := float64(len(stages) - 1)
	for index, stage := range stages {
		logging.Infof("Uninstall stage %s", stage.Name())
		p.progressBar.SetValue(float64(index) / total)
		p.statusLabel.SetText(stage.Name())
		err := stage.Execute()
		logging.LogError(err)
		if err != nil {
			line := fmt.Sprintf("\n%s: %v", stage.Name(), err.Error())
			p.errorsLabel.SetText(p.errorsLabel.Text + line)
		}
	}
}
