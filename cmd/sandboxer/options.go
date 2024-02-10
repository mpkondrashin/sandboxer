package main

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/vone"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
)

type OptionsWindow struct {
	ModalWindow
	conf         *config.Configuration
	tokenEntry   *widget.Entry
	domainLabel  *widget.Label
	cancelDetect context.CancelFunc
}

func NewOptionsWindow(modalWindow ModalWindow, conf *config.Configuration) *OptionsWindow {
	s := &OptionsWindow{
		ModalWindow: modalWindow,
		conf:        conf,
	}

	labelTop := widget.NewLabel("Please open Vision One console to get all nessesary parameters")
	s.tokenEntry = widget.NewMultiLineEntry()
	s.tokenEntry.Wrapping = fyne.TextWrapBreak
	s.tokenEntry.OnChanged = s.DetectDomain
	tokenFormItem := widget.NewFormItem("Token:", s.tokenEntry)
	tokenFormItem.HintText = "Go to Administrator -> API Keys"
	//roleHint := "Go to Administration -> User Roles -> Permissions -> Threat Intelligence -> Sandbox Analysis -> \"View, filter, and search\" and\"Submit object\""
	// apiKeyHitn := "Go to Administration -> API Keys -> Add API Key"
	s.domainLabel = widget.NewLabel("")
	domainFormItem := widget.NewFormItem("Domain:", s.domainLabel)

	optionsForm := widget.NewForm(
		tokenFormItem,
		domainFormItem,
	)

	saveButton := widget.NewButton("Save", s.Save)
	cancelButton := widget.NewButton("Cancel", s.Hide)
	bottons := container.NewHBox(cancelButton, saveButton)
	// add link to open v1 console(?)
	s.win.SetContent(container.NewVBox(labelTop, optionsForm, bottons))
	return s
}

func (s *OptionsWindow) Save() {
	s.conf.VisionOne.Token = s.tokenEntry.Text
	if s.domainLabel.Text != ErrorDomain {
		s.conf.VisionOne.Domain = s.domainLabel.Text
	}
	if err := s.conf.Save(); err != nil {
		logging.Errorf("Save Config: %v", err)
		dialog.ShowError(err, s.win)
		return
	}
	s.Hide()
}

/*
	func (s *OptionsWindow) Cancel() {
		s.Hide()
	}
*/
const ErrorDomain = "Error"

func (s *OptionsWindow) Update() {
	s.tokenEntry.SetText(s.conf.VisionOne.Token)
	if s.conf.VisionOne.Domain == "" {
		s.conf.VisionOne.Domain = vone.DetectVisionOneDomain(context.TODO(), s.conf.VisionOne.Token)
	}
	if s.conf.VisionOne.Domain != "" {
		s.domainLabel.SetText(s.conf.VisionOne.Domain)
	} else {
		s.domainLabel.SetText(ErrorDomain)
	}
}

func (s *OptionsWindow) DetectDomain(token string) {
	go func() {
		if s.cancelDetect != nil {
			log.Println("not nil - cancel")
			s.cancelDetect()
		}
		var ctx context.Context
		ctx, s.cancelDetect = context.WithCancel(context.TODO())
		defer func() {
			log.Println("defer cancel")
			if s.cancelDetect != nil {
				s.cancelDetect()
			}
			s.cancelDetect = nil
		}()
		domain := vone.DetectVisionOneDomain(ctx, token)
		if domain != "" {
			s.domainLabel.SetText(domain)
		} else {
			s.domainLabel.SetText(ErrorDomain)
		}
	}()
}

func (s *OptionsWindow) Show() {
	s.win.Show()
	s.Update()
}
