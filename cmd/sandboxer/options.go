/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

options.go

Options window
*/
package main

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/ddan"
	"github.com/mpkondrashin/vone"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
)

type OptionsWindow struct {
	conf *config.Configuration

	voneCheck    *widget.Check
	tokenEntry   *widget.Entry
	domainLabel  *widget.Label
	cancelDetect context.CancelFunc

	ddanCheck          *widget.Check
	ddanURLEntry       *widget.Entry
	ddanAPIKeyEntry    *widget.Entry
	ddanIgnoreTLSCheck *widget.Check
	ddanTest           *widget.Label
	cancelTestDDAn     context.CancelFunc

	ignoreEntry   *widget.Entry
	tasksKeepDays *widget.Entry
}

func NewOptionsWindow(conf *config.Configuration) *OptionsWindow {
	return &OptionsWindow{
		conf: conf,
	}
}

func (s *OptionsWindow) Show() {
	//s.Update()
}

func (s *OptionsWindow) Hide() {}

func (s *OptionsWindow) Name() string {
	return "Options"
}

func (s *OptionsWindow) Content(w *ModalWindow) fyne.CanvasObject {

	saveButton := widget.NewButton("Save", func() { s.Save(w) })
	cancelButton := widget.NewButton("Cancel", w.Hide)
	buttons := container.NewHBox(cancelButton, saveButton)
	// add link to open v1 console(?)

	ddanTab := container.NewTabItem("Analyzer", s.DDAnSettings())
	voneTab := container.NewTabItem("Vision One", s.VisionOneSettings())
	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", s.GeneralSettings()),
		voneTab,
		ddanTab,
	)
	tabs.OnSelected = func(tab *container.TabItem) {
		switch tab {
		case ddanTab:
			s.TestAnalyzer()
		case voneTab:
			s.Update()
		}
	}
	return container.NewVBox(tabs, buttons)
}

func (s *OptionsWindow) DDAnSettings() fyne.CanvasObject {
	labelTop := widget.NewLabel("Deep Discovery Analyzer settings")

	s.ddanCheck = widget.NewCheck("Use Analyzer sandbox", func(checked bool) {
		s.voneCheck.Checked = !checked
	})
	s.ddanCheck.Checked = s.conf.SandboxType == config.SandboxAnalyzer

	s.ddanURLEntry = widget.NewEntry()
	s.ddanURLEntry.SetText(s.conf.DDAn.URL)
	s.ddanURLEntry.OnChanged = func(string) {
		s.TestAnalyzer()
	}
	urlFormItem := widget.NewFormItem("Address:", s.ddanURLEntry)
	urlFormItem.HintText = "DNS name or IP address"

	s.ddanAPIKeyEntry = widget.NewEntry()
	s.ddanAPIKeyEntry.SetText(s.conf.DDAn.APIKey)
	s.ddanAPIKeyEntry.OnChanged = func(string) {
		s.TestAnalyzer()
	}
	apiKeyFormItem := widget.NewFormItem("API Key:", s.ddanAPIKeyEntry)
	apiKeyFormItem.HintText = "Go to Help -> About on Analyzer console"

	s.ddanIgnoreTLSCheck = widget.NewCheck("Ignore", nil)
	s.ddanIgnoreTLSCheck.SetChecked(s.conf.DDAn.IgnoreTLSErrors)
	s.ddanIgnoreTLSCheck.OnChanged = func(bool) {
		s.TestAnalyzer()
	}
	ignoreTLSFormItem := widget.NewFormItem("TLS Errors: ", s.ddanIgnoreTLSCheck)

	s.ddanTest = widget.NewLabel("")

	ddanForm := widget.NewForm(urlFormItem, apiKeyFormItem, ignoreTLSFormItem)
	return container.NewVBox(s.ddanCheck, labelTop, ddanForm, s.ddanTest)
}

func (s *OptionsWindow) VisionOneSettings() fyne.CanvasObject {
	labelTop := widget.NewLabel("Please open Vision One console to get all nessesary parameters")

	s.voneCheck = widget.NewCheck("Use Vision One sandbox", func(checked bool) {
		s.ddanCheck.Checked = !checked
	})
	s.voneCheck.Checked = s.conf.SandboxType == config.SandboxVisionOne

	s.tokenEntry = widget.NewMultiLineEntry()
	s.tokenEntry.SetText(s.conf.VisionOne.Token)
	s.tokenEntry.Wrapping = fyne.TextWrapBreak
	s.tokenEntry.OnChanged = s.DetectDomain
	tokenFormItem := widget.NewFormItem("Token:", s.tokenEntry)
	tokenFormItem.HintText = "Go to Administrator -> API Keys"
	//roleHint := "Go to Administration -> User Roles -> Permissions -> Threat Intelligence -> Sandbox Analysis -> \"View, filter, and search\" and\"Submit object\""
	// apiKeyHitn := "Go to Administration -> API Keys -> Add API Key"

	s.domainLabel = widget.NewLabel(s.conf.VisionOne.Domain)
	domainFormItem := widget.NewFormItem("Domain:", s.domainLabel)
	optionsForm := widget.NewForm(
		tokenFormItem,
		domainFormItem,
	)
	return container.NewVBox(s.voneCheck, labelTop, optionsForm)
}

func (s *OptionsWindow) GeneralSettings() fyne.CanvasObject {
	settingsLabel := widget.NewLabel("General Options")
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
	return container.NewVBox(settingsLabel, settingsForm)
}

func (s *OptionsWindow) Save(w *ModalWindow) {

	if s.ddanCheck.Checked {
		s.conf.SandboxType = config.SandboxAnalyzer
	}
	if s.voneCheck.Checked {
		s.conf.SandboxType = config.SandboxVisionOne
	}

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

	s.conf.DDAn.URL = s.GetDDAnURL()
	s.conf.DDAn.APIKey = s.ddanAPIKeyEntry.Text
	s.conf.DDAn.IgnoreTLSErrors = s.ddanIgnoreTLSCheck.Checked

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

func (s *OptionsWindow) TestAnalyzer() {
	log.Println("TestAnalyzer")
	go func() {
		log.Println("TestAnalyzer go")
		if s.cancelTestDDAn != nil {
			log.Println("s.cancelTestDDAn != nil ")
			s.cancelTestDDAn()
		}
		var ctx context.Context
		ctx, s.cancelTestDDAn = context.WithCancel(context.TODO())
		defer func() {
			if s.cancelTestDDAn != nil {
				s.cancelTestDDAn()
			}
			s.cancelTestDDAn = nil
		}()
		s.ddanTest.SetText("Checking connection...")
		u, err := url.Parse(s.GetDDAnURL())
		if err != nil {
			s.ddanTest.SetText(err.Error())
			return
		}
		log.Println("TestAnalyzer u = ", u)
		apiKey := strings.TrimSpace(s.ddanAPIKeyEntry.Text)
		log.Println("apiKey = ", apiKey)
		analyzer := ddan.NewClient(s.conf.DDAn.ProductName, s.conf.DDAn.Hostname).
			SetAnalyzer(u, apiKey, s.ddanIgnoreTLSCheck.Checked)
		log.Println("analyzer ", analyzer)
		if s.conf.DDAn.ProtocolVersion != "" {
			log.Println("analyzer set version ", s.conf.DDAn.ProtocolVersion)
			analyzer.SetProtocolVersion(s.conf.DDAn.ProtocolVersion)
		}
		log.Println("To test connection")
		ctxTimeout, cancelTimeout := context.WithTimeout(ctx, 5*time.Second)
		defer cancelTimeout()
		err = analyzer.TestConnection(ctxTimeout)
		log.Println("TestConnection err ", err)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				if errors.Is(err, context.DeadlineExceeded) {
					s.ddanTest.SetText("Connection timed out")
				} else {
					s.ddanTest.SetText(err.Error())
				}
			}
		} else {
			s.ddanTest.SetText("Connection is Ok")
		}
	}()
}

func (s *OptionsWindow) GetDDAnURL() (result string) {
	result = strings.TrimSpace(s.ddanURLEntry.Text)
	if strings.HasPrefix(result, "https://") {
		return
	}
	if strings.HasPrefix(result, "http://") {
		return
	}
	return "https://" + result
}
