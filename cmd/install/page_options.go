package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageOptions struct {
	//regionList     *widget.Select
	//accountIDEntry *widget.Entry
	tokenEntry *widget.Entry
	//awsRegionList  *widget.Select
}

var _ Page = &PageOptions{}

func (p *PageOptions) Name() string {
	return "Options"
}

func (p *PageOptions) Content(win fyne.Window, model *Model) fyne.CanvasObject {
	labelTop := widget.NewLabel("Please open Vision One console to get all nessesary parameters")
	//p.accountIDEntry = widget.NewEntry()
	//p.accountIDEntry.Text = model.config.AccountID
	//p.accountIDEntry.Validator = ValidateAccountID
	p.tokenEntry = widget.NewMultiLineEntry()
	p.tokenEntry.Text = model.config.Token
	p.tokenEntry.Wrapping = fyne.TextWrapBreak
	tokenFormItem := widget.NewFormItem("Token:", p.tokenEntry)
	tokenFormItem.HintText = "Go to XXXXXXX"

	optionsForm := widget.NewForm(
		tokenFormItem,
	)
	return container.NewVBox(labelTop, optionsForm)
}

func (p *PageOptions) AquireData(model *Model) error {
	model.config.Token = p.tokenEntry.Text
	//model.config.AccountID = p.accountIDEntry.Text

	return nil
}
