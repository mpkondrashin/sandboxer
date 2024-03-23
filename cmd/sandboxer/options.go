/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

options.go

Options window
*/
package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/settings"
)

type OptionsWindow struct {
	conf *config.Configuration

	voneCheck    *widget.Check
	voneSettings *settings.VisionOne

	ddanCheck    *widget.Check
	ddanSettings *settings.DDAn

	proxySettings *settings.Proxy

	ignoreEntry       *widget.Entry
	tasksKeepDays     *widget.Entry
	showNotifications *widget.Check
}

func NewOptionsWindow(conf *config.Configuration) *OptionsWindow {
	return &OptionsWindow{
		conf:          conf,
		voneSettings:  settings.NewVisionOne(conf.VisionOne),
		ddanSettings:  settings.NewDDAnSettings(conf.DDAn),
		proxySettings: settings.NewProxy(conf.Proxy),
	}
}

func (s *OptionsWindow) Show() {}

func (s *OptionsWindow) Hide() {}

func (s *OptionsWindow) Name() string {
	return "Options"
}

func (s *OptionsWindow) Icon() fyne.Resource {
	return theme.SettingsIcon()
}

func (s *OptionsWindow) Content(w *ModalWindow) fyne.CanvasObject {
	saveButton := widget.NewButton("Save", func() { s.Save(w) })
	cancelButton := widget.NewButton("Cancel", w.Hide)
	buttons := container.NewHBox(cancelButton, saveButton)
	// add link to open v1 console(?)

	ddanTab := container.NewTabItem("Analyzer", s.DDAnSettings(w))
	voneTab := container.NewTabItem("Vision One", s.VisionOneSettings())
	proxyTab := container.NewTabItem("Proxy", s.ProxySettings())
	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", s.GeneralSettings()),
		voneTab,
		ddanTab,
		proxyTab,
	)
	tabs.OnSelected = func(tab *container.TabItem) {
		switch tab {
		case ddanTab:
			s.ddanSettings.Update()
		case voneTab:
			s.voneSettings.Update()
		}
	}
	return container.NewVBox(tabs, buttons)
}

func (s *OptionsWindow) DDAnSettings(w *ModalWindow) fyne.CanvasObject {
	labelTop := widget.NewLabel("Deep Discovery Analyzer settings")

	s.ddanCheck = widget.NewCheck("Use Analyzer sandbox", func(checked bool) {
		s.voneCheck.Checked = !checked
	})
	s.ddanCheck.Checked = s.conf.SandboxType == config.SandboxAnalyzer
	unregisterButton := widget.NewButton("Unregister", func() {
		logging.Infof("Unregister from Analyzer")
		analyzer, err := s.conf.DDAn.Analyzer()
		if err != nil {
			logging.LogError(err)
			dialog.ShowError(err, w.win)
			return
		}
		if err := analyzer.Unregister(context.TODO()); err != nil {
			logging.LogError(err)
			dialog.ShowError(err, w.win)
			return
		}
		dialog.ShowInformation("Analyzer", "Successfully unregistered", w.win)
	})
	return container.NewVBox(s.ddanCheck, labelTop, s.ddanSettings.Widget(), unregisterButton)
}

func (s *OptionsWindow) VisionOneSettings() fyne.CanvasObject {
	labelTop := widget.NewLabel("Please open Vision One console to get all nessesary parameters")

	s.voneCheck = widget.NewCheck("Use Vision One sandbox", func(checked bool) {
		s.ddanCheck.Checked = !checked
	})
	s.voneCheck.Checked = s.conf.SandboxType == config.SandboxVisionOne

	return container.NewVBox(s.voneCheck, labelTop, s.voneSettings.Widget())
}

func (s *OptionsWindow) ProxySettings() fyne.CanvasObject {
	labelTop := widget.NewLabel("Proxy settings")

	return container.NewVBox(labelTop, s.proxySettings.Widget())
}

func (s *OptionsWindow) GeneralSettings() fyne.CanvasObject {
	settingsLabel := widget.NewLabel("General Options")
	s.ignoreEntry = widget.NewEntry()
	s.ignoreEntry.SetText(strings.Join(s.conf.Ignore, ", "))
	ignoreFormItem := widget.NewFormItem("Ignore:", s.ignoreEntry)
	ignoreFormItem.HintText = "Comma-separated list of file masks"

	s.tasksKeepDays = widget.NewEntry()
	s.tasksKeepDays.SetText(strconv.Itoa(s.conf.GetTasksKeepDays()))
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

	s.showNotifications = widget.NewCheck("Show", nil)
	s.showNotifications.Checked = s.conf.GetShowNotifications()
	notificatonsFormItem := widget.NewFormItem("Notifications:", s.showNotifications)

	settingsForm := widget.NewForm(ignoreFormItem, tasksKeepDaysFormItem, notificatonsFormItem)
	return container.NewVBox(settingsLabel, settingsForm)
}

func (s *OptionsWindow) Save(w *ModalWindow) {

	s.conf.Ignore = nil
	for _, ign := range strings.Split(s.ignoreEntry.Text, ",") {
		ign := strings.TrimSpace(ign)
		if len(ign) > 0 {
			s.conf.Ignore = append(s.conf.Ignore, ign)
		}
	}
	days, err := strconv.Atoi(s.tasksKeepDays.Text)
	if err == nil {
		s.conf.SetTasksKeepDays(days)
	}
	s.conf.SetShowNotifications(s.showNotifications.Checked)

	if s.ddanCheck.Checked {
		s.conf.SandboxType = config.SandboxAnalyzer
	}
	if s.voneCheck.Checked {
		s.conf.SandboxType = config.SandboxVisionOne
	}

	if err := s.voneSettings.Aquire(); err != nil {
		err = fmt.Errorf("Vision One Setting: %v", err)
		logging.LogError(err)
		dialog.ShowError(err, w.win)
		return
	}

	if err = s.ddanSettings.Aquire(); err != nil {
		err = fmt.Errorf("Analyzer Settings: %v", err)
		logging.LogError(err)
		dialog.ShowError(err, w.win)
		return
	}

	if err := s.proxySettings.Aquire(); err != nil {
		err = fmt.Errorf("Proxy Setting: %v", err)
		logging.LogError(err)
		dialog.ShowError(err, w.win)
		return
	}

	if err := s.conf.Save(); err != nil {
		logging.Errorf("Save Config: %v", err)
		dialog.ShowError(err, w.win)
		return
	}
	w.Hide()
}
