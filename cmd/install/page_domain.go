/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_domain.go

Pick Vision One domain
*/
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
	chooseLabel      *widget.Label
	visionOneDomains *widget.Select
}

var _ Page = &PageDomain{}

func (p *PageDomain) Name() string {
	return "Domain"
}

/*
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

		p.chooseLabel = widget.NewLabel(label)
		var domains []string
		for _, d := range vone.RegionalDomains {
			//l := fmt.Sprintf("%s (%s)", d.Region, d.Domain)
			domains = append(domains, d.Domain)
		}

		p.visionOneDomains = widget.NewSelect(domains, nil)
		if selectedDomain != "" {
			p.visionOneDomains.SetSelected(selectedDomain)
		}
		domainForm := widget.NewForm(
			widget.NewFormItem("Domain:", p.visionOneDomains),
		)
		return container.NewVBox(p.chooseLabel, domainForm)
	}
*/
func (p *PageDomain) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {

	p.chooseLabel = widget.NewLabel("Choose Vision One Domain:")
	var domains []string
	for _, d := range vone.RegionalDomains {
		//l := fmt.Sprintf("%s (%s)", d.Region, d.Domain)
		domains = append(domains, d.Domain)
	}

	p.visionOneDomains = widget.NewSelect(domains, nil)
	if installer.config.VisionOne.Domain != "" {
		p.visionOneDomains.SetSelected(installer.config.VisionOne.Domain)
	}
	domainForm := widget.NewForm(
		widget.NewFormItem("Domain:", p.visionOneDomains),
	)
	return container.NewVBox(p.chooseLabel, domainForm)
}
func (p *PageDomain) Run(win fyne.Window, installer *Installer) {
	if p.visionOneDomains.Selected != "" {
		return
	}
	detected := vone.DetectVisionOneDomain(context.TODO(), installer.config.VisionOne.Token)
	if detected == "" {
		return
	}
	p.visionOneDomains.SetSelected(detected)
	p.chooseLabel.SetText("Detected Vision One Domain:")
}

func (p *PageDomain) AquireData(installer *Installer) error {
	if p.visionOneDomains.Selected == "" {
		return fmt.Errorf("No Domain selected")
	}
	installer.config.VisionOne.Domain = p.visionOneDomains.Selected
	return nil
}
