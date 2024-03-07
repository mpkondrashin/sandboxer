/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

quota.go

Quota window
*/
package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
)

type QuotaWindow struct {
	win        fyne.Window
	conf       *config.Configuration
	reserve    binding.String
	submission binding.String
	exemption  binding.String
	remaining  binding.String
}

func NewQuotaWindow(conf *config.Configuration) *QuotaWindow {
	w := &QuotaWindow{
		conf:       conf,
		reserve:    binding.NewString(),
		submission: binding.NewString(),
		exemption:  binding.NewString(),
		remaining:  binding.NewString(),
	}
	return w
}

func (w *QuotaWindow) Content(modal *ModalWindow) fyne.CanvasObject {
	w.win = modal.win
	reserveCountItem := widget.NewFormItem("Daily Reserve:", widget.NewLabelWithData(w.reserve))
	submissionCountItem := widget.NewFormItem("Files Submitted:", widget.NewLabelWithData(w.submission))
	exemptionCountItem := widget.NewFormItem("Unsupported Files Submitted:", widget.NewLabelWithData(w.exemption))
	remainingCountItem := widget.NewFormItem("Remaining:", widget.NewLabelWithData(w.remaining))
	updateItem := widget.NewFormItem("", widget.NewButton("Update", func() {
		w.Update()
	}))
	form := widget.NewForm(
		reserveCountItem,
		submissionCountItem,
		exemptionCountItem,
		remainingCountItem,
		updateItem,
	)
	w.Reset()
	return form

}

func (w *QuotaWindow) Name() string {
	return "Vision One Quota"
}

func (s *QuotaWindow) Reset() {
	s.reserve.Set("?")
	s.submission.Set("?")
	s.exemption.Set("?")
	s.remaining.Set("?")
}

func (s *QuotaWindow) Update() {
	s.Reset()
	vOne, err := s.conf.VisionOne.VisionOneSandbox()
	if err != nil {
		logging.LogError(err)
		dialog.ShowError(err, s.win)
		return
	}
	result, err := vOne.SandboxDailyReserve().Do(context.TODO())
	if err != nil {
		logging.LogError(err)
		dialog.ShowError(err, s.win)
		return
	}
	reserveCount := fmt.Sprintf("%d", result.SubmissionReserveCount)
	submissionCount := fmt.Sprintf("%d", result.SubmissionCount)
	exemptionCount := fmt.Sprintf("%d", result.SubmissionExemptionCount)
	remainingCount := fmt.Sprintf("%d                            ", result.SubmissionRemainingCount)

	s.reserve.Set(reserveCount)
	s.submission.Set(submissionCount)
	s.exemption.Set(exemptionCount)
	s.remaining.Set(remainingCount)
}

func (s *QuotaWindow) Show() {
	go s.Update()
}

func (s *QuotaWindow) Hide() {}
