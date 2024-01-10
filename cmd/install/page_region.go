package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mpkondrashin/vone"
)

type PageDomain struct {
	//regionList     *widget.Select
	visionOneDomains *widget.Select
}

var _ Page = &PageDomain{}

func (p *PageDomain) Name() string {
	return "Domain"
}

func (p *PageDomain) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {

	selectedDomain := ""
	label := "Choose Vision One Domain:"
	if installer.config.Domain != "" {
		selectedDomain = installer.config.Domain
	} else {
		selectedDomain = vone.DetectVisionOneDomain(context.TODO(), installer.config.Token)
		if selectedDomain != "" {
			label = "Detected Vision One Domain:"
		}
	}

	labelTop := widget.NewLabel(label)
	var domains []string
	for _, d := range vone.RegionalDomains {
		//l := fmt.Sprintf("%s (%s)", d.Region, d.Domain)
		domains = append(domains, d.Domain)
	}

	p.visionOneDomains = widget.NewSelect(domains, nil)
	if selectedDomain != "" {
		p.visionOneDomains.SetSelected(selectedDomain)
	}
	passwordForm := widget.NewForm(
		widget.NewFormItem("Domain:", p.visionOneDomains),
	)
	return container.NewVBox(labelTop, passwordForm)
}

func (p *PageDomain) AquireData(installer *Installer) error {
	if p.visionOneDomains.Selected == "" {
		return fmt.Errorf("No Domain selected")
	}
	installer.config.Domain = p.visionOneDomains.Selected
	return nil
}
