/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

options.go

Options window
*/
package main

import (
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	ignoreEntry   *widget.Entry
	tasksKeepDays *widget.Entry
}

func NewOptionsWindow(conf *config.Configuration) *OptionsWindow {
	return &OptionsWindow{
		conf:         conf,
		voneSettings: settings.NewVisionOne(&conf.VisionOne),
		ddanSettings: settings.NewDDAnSettings(&conf.DDAn),
	}
}

func (s *OptionsWindow) Show() {}

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
			s.ddanSettings.Update()
		case voneTab:
			s.voneSettings.Update()
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

	return container.NewVBox(s.ddanCheck, labelTop, s.ddanSettings.Widget())
}

func (s *OptionsWindow) VisionOneSettings() fyne.CanvasObject {
	labelTop := widget.NewLabel("Please open Vision One console to get all nessesary parameters")

	s.voneCheck = widget.NewCheck("Use Vision One sandbox", func(checked bool) {
		s.ddanCheck.Checked = !checked
	})
	s.voneCheck.Checked = s.conf.SandboxType == config.SandboxVisionOne

	return container.NewVBox(s.voneCheck, labelTop, s.voneSettings.Widget())
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

	s.voneSettings.Aquire()
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

	s.ddanSettings.Aquire()

	if err := s.conf.Save(); err != nil {
		logging.Errorf("Save Config: %v", err)
		dialog.ShowError(err, w.win)
		return
	}
	w.Hide()
}
