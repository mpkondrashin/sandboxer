package main

import (
	"context"
	"examen/pkg/config"
	"examen/pkg/logging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/vone"
)

type OptionsWindow struct {
	conf        *config.Configuration
	win         fyne.Window
	tokenEntry  *widget.Entry
	domainLabel *widget.Label
}

func NewOptionsWindow(app fyne.App, conf *config.Configuration) *OptionsWindow {
	s := &OptionsWindow{
		conf: conf,
		win:  app.NewWindow("Options"),
	}
	s.win.SetCloseIntercept(func() {
		s.win.Hide()
	})

	labelTop := widget.NewLabel("Please open Vision One console to get all nessesary parameters")
	s.tokenEntry = widget.NewMultiLineEntry()
	s.tokenEntry.Wrapping = fyne.TextWrapBreak
	s.tokenEntry.OnChanged = s.DetectDomain
	tokenFormItem := widget.NewFormItem("Token:", s.tokenEntry)
	tokenFormItem.HintText = "Go to XXXXXXX"

	s.domainLabel = widget.NewLabel("")
	domainFormItem := widget.NewFormItem("Domain:", s.domainLabel)

	optionsForm := widget.NewForm(
		tokenFormItem,
		domainFormItem,
	)

	saveButton := widget.NewButton("Save", s.Save)
	cancelButton := widget.NewButton("Cancel", s.Cancel)
	bottons := container.NewHBox(cancelButton, saveButton)
	// add link to open v1 console(?)
	s.win.SetContent(container.NewVBox(labelTop, optionsForm, bottons))
	return s
}

func (s *OptionsWindow) Save() {
	//	conf, err := config.LoadConfiguration(globals.AppID, globals.ConfigFileName)
	//	if err != nil {
	//		logging.Errorf("LoadConfig: %v", err)
	//		dialog.ShowError(err, s.win)
	//		return
	//	}
	s.conf.VisionOne.Token = s.tokenEntry.Text
	if err := s.conf.Save(); err != nil {
		logging.Errorf("Save Config: %v", err)
		dialog.ShowError(err, s.win)
		return
	}
	s.win.Hide()
}

func (s *OptionsWindow) Cancel() {
	s.win.Hide()
}

func (s *OptionsWindow) Update() {
	s.tokenEntry.SetText(s.conf.VisionOne.Token) // ???
	if s.conf.VisionOne.Domain == "" {
		s.conf.VisionOne.Domain = vone.DetectVisionOneDomain(context.TODO(), s.conf.VisionOne.Token)
	}
	s.domainLabel.SetText(s.conf.VisionOne.Domain)
}

func (s *OptionsWindow) DetectDomain(token string) {
	domain := vone.DetectVisionOneDomain(context.TODO(), token)
	s.domainLabel.SetText(domain)
}

func (s *OptionsWindow) Show() {
	s.win.Show()
	s.Update()
}
