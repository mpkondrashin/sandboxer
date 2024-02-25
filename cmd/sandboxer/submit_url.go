/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

submit_url.go

Submit URL window
*/
package main

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
)

type SubmitURLWindow struct {
	list         *task.TaskList
	channels     *task.Channels
	urlEntry     *widget.Entry
	submitButton *widget.Button
}

func NewSubmitURLWindow(list *task.TaskList, channels *task.Channels) *SubmitURLWindow {
	return &SubmitURLWindow{
		list:     list,
		channels: channels,
	}
}

func (s *SubmitURLWindow) Show() {
	s.urlEntry.SetText("")
}

func (s *SubmitURLWindow) Hide() {}

func (s *SubmitURLWindow) Name() string {
	return "Submit URL"
}

func (s *SubmitURLWindow) Content(w *ModalWindow) fyne.CanvasObject {
	w.win.Resize(fyne.Size{Width: 400, Height: 40})
	labelTop := widget.NewLabel("Enter URL")
	s.urlEntry = widget.NewEntry()
	s.urlEntry.OnChanged = s.Update
	tokenFormItem := widget.NewFormItem("URL:", s.urlEntry)
	optionsForm := widget.NewForm(
		tokenFormItem,
	)
	s.submitButton = widget.NewButton("Submit", func() {
		s.Submit()
		w.Hide()
	})
	s.submitButton.Disable()
	cancelButton := widget.NewButton("Cancel", w.Hide)
	bottons := container.NewHBox(cancelButton, s.submitButton)
	// add link to open v1 console(?)
	return container.NewVBox(labelTop, optionsForm, bottons)
}

func (s *SubmitURLWindow) Submit() {
	url := strings.TrimSpace(s.urlEntry.Text)
	tsk, err := s.list.NewTask(task.URLTask, url)
	if err != nil {
		if !errors.Is(err, task.ErrAlreadyExists) {
			logging.LogError(err)
		}
		return
	}
	s.channels.TaskChannel[task.ChPrefilter] <- tsk
}

func (s *SubmitURLWindow) Update(str string) {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		s.submitButton.Disable()
	} else {
		s.submitButton.Enable()
	}
}

//var urlRegex = regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)

func IsURL(s string) bool {
	return true
	//_, err := url.ParseRequestURI(s)
	//return err != nil
	//return govalidator.IsRequestURL(s)
	/*
		ts := strings.TrimPrefix(s, "http://")
		ts = strings.TrimPrefix(ts, "https://")
		if len(ts) == len(s) {
			return false
		}

		return urlRegex.MatchString(s)*/
}
