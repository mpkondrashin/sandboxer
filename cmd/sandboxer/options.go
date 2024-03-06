/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

options.go

Options window
*/
package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/vone"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
)

type OptionsWindow struct {
	conf          *config.Configuration
	tokenEntry    *widget.Entry
	domainLabel   *widget.Label
	ignoreEntry   *widget.Entry
	tasksKeepDays *widget.Entry
	cancelDetect  context.CancelFunc
}

func NewOptionsWindow(conf *config.Configuration) *OptionsWindow {
	return &OptionsWindow{
		conf: conf,
	}
}

func (s *OptionsWindow) Show() {
	s.Update()
}

func (s *OptionsWindow) Hide() {}

func (s *OptionsWindow) Name() string {
	return "Options"
}

func (s *OptionsWindow) Content(w *ModalWindow) fyne.CanvasObject {
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

	settingsLabel := widget.NewLabel("Settings")
	s.ignoreEntry = widget.NewEntry()
	s.ignoreEntry.SetText(strings.Join(s.conf.Ignore, ", "))
	ignoreFormItem := widget.NewFormItem("Ignore:", s.ignoreEntry)
	ignoreFormItem.HintText = "Comma-separated list of file masks"

	s.tasksKeepDays = widget.NewEntry()
	s.tasksKeepDays.SetText(strconv.Itoa(s.conf.TasksKeepDays))
	s.tasksKeepDays.OnChanged = func(str string) {
		n := ""
		for _, ch := range str {
			if unicode.IsDigit(ch) {
				n += string(ch)
			}
		}
		if n != str {
			s.tasksKeepDays.SetText(n)
		}
	}
	tasksKeepDaysFormItem := widget.NewFormItem("Delete tasks after: ", s.tasksKeepDays)
	tasksKeepDaysFormItem.HintText = "Number of days"

	settingsForm := widget.NewForm(ignoreFormItem, tasksKeepDaysFormItem)
	settingsVBox := container.NewVBox(settingsLabel, settingsForm)

	saveButton := widget.NewButton("Save", func() { s.Save(w) })
	cancelButton := widget.NewButton("Cancel", w.Hide)
	buttons := container.NewHBox(cancelButton, saveButton)
	// add link to open v1 console(?)
	vOneVBox := container.NewVBox(labelTop, optionsForm)

	tabs := container.NewAppTabs(
		container.NewTabItem("Vision One", vOneVBox),
		container.NewTabItem("Setting", settingsVBox),
	)

	return container.NewVBox(tabs, buttons)
}

func (s *OptionsWindow) Save(w *ModalWindow) {
	s.conf.VisionOne.Token = strings.TrimSpace(s.tokenEntry.Text)
	if s.domainLabel.Text != ErrorDomain {
		s.conf.VisionOne.Domain = s.domainLabel.Text
	}
	s.conf.Ignore = nil
	for _, ign := range strings.Split(s.ignoreEntry.Text, ",") {
		ign := strings.TrimSpace(ign)
		if len(ign) > 0 {
			s.conf.Ignore = append(s.conf.Ignore, ign)
		}
	}
	days, err := strconv.Atoi(s.tasksKeepDays.Text)
	if err == nil {
		s.conf.TasksKeepDays = days
	}
	if err := s.conf.Save(); err != nil {
		logging.Errorf("Save Config: %v", err)
		dialog.ShowError(err, w.win)
		return
	}
	w.Hide()
}

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
