package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PagePassword struct {
	passwordEntry *widget.Entry
}

var _ Page = &PagePassword{}

func (p *PagePassword) Name() string {
	return "Auth"
}

func (p *PagePassword) Content(win fyne.Window, model *Model) fyne.CanvasObject {
	configExist, err := model.ConfigExists()
	if err != nil {
		dialog.ShowError(err, win)
	}
	labelText := "Please provide password that will be used to encrypt API key"
	if configExist {
		labelText = "Please provide password to decrypt API key"
	}
	labelTop := widget.NewLabel(labelText)
	p.passwordEntry = widget.NewPasswordEntry()
	p.passwordEntry.Text = model.password
	p.passwordEntry.Validator = CheckPassword
	passwordFormItem := widget.NewFormItem("Password:", p.passwordEntry)
	if !configExist {
		passwordFormItem.HintText = "At least 8 characters, upper/lower case, digits and special characters"
	}
	passwordForm := widget.NewForm(passwordFormItem)
	return container.NewVBox(labelTop, passwordForm)
}

func (p *PagePassword) AquireData(model *Model) error {
	if err := p.passwordEntry.Validate(); err != nil {
		return err
	}
	model.password = p.passwordEntry.Text
	return model.LoadConfig()
}
