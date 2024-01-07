package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mpkondrashin/vone"
)

type PageDomain struct {
	//regionList     *widget.Select
	awsRegionList *widget.Select
}

var _ Page = &PageDomain{}

func (p *PageDomain) Name() string {
	return "Domain"
}

func (p *PageDomain) Content(win fyne.Window, model *Model) fyne.CanvasObject {

	selectedRegion := ""
	if model.config.Domain != "" {
		selectedRegion = model.config.Domain
	} else {
	}

	labelTop := widget.NewLabel("Choose Vision One Domain")
	var domains []string
	for _, d := range vone.RegionalDomains {
		l := fmt.Sprintf("%s (%s)", d.Region, d.Domain)
		domains = append(domains, l)
	}

	p.awsRegionList = widget.NewSelect(domains, nil)
	if selectedRegion != "" {
		p.awsRegionList.SetSelected(selectedRegion)
	}
	passwordForm := widget.NewForm(
		widget.NewFormItem("Region:", p.awsRegionList),
	)
	return container.NewVBox(labelTop, passwordForm)
}

func (p *PageDomain) AquireData(model *Model) error {
	model.config.Domain = p.awsRegionList.Selected
	if model.Changed() {
		return model.Save()
	}
	return nil
}
