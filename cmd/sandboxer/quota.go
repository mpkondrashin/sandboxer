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

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/mpkondrashin/vone"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
)

type QuotaWindow struct {
	ModalWindow
	conf       *config.Configuration
	reserve    binding.String
	submission binding.String
	exemption  binding.String
	remaining  binding.String
}

func NewQuotaWindow(modalWindow ModalWindow, conf *config.Configuration) *QuotaWindow {
	s := &QuotaWindow{
		ModalWindow: modalWindow,
		conf:        conf,
		reserve:     binding.NewString(),
		submission:  binding.NewString(),
		exemption:   binding.NewString(),
		remaining:   binding.NewString(),
	}
	s.Reset()
	reserveCountItem := widget.NewFormItem("Daily Reserve:", widget.NewLabelWithData(s.reserve))
	submissionCountItem := widget.NewFormItem("Files Submitted:", widget.NewLabelWithData(s.submission))
	exemptionCountItem := widget.NewFormItem("Unsupported Files Submitted:", widget.NewLabelWithData(s.exemption))
	remainingCountItem := widget.NewFormItem("Remaining:", widget.NewLabelWithData(s.remaining))
	updateItem := widget.NewFormItem("", widget.NewButton("Update", func() {
		s.Update()
	}))
	form := widget.NewForm(
		reserveCountItem,
		submissionCountItem,
		exemptionCountItem,
		remainingCountItem,
		updateItem,
	)
	s.win.SetContent(form)
	return s
}

func (s *QuotaWindow) Reset() {
	s.reserve.Set("?")
	s.submission.Set("?")
	s.exemption.Set("?")
	s.remaining.Set("?")
}

func (s *QuotaWindow) Update() {
	s.Reset()
	vOne := vone.NewVOne(s.conf.VisionOne.Domain, s.conf.VisionOne.Token)
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
	//exemptionCount := fmt.Sprintf("%d (Number of samples submitted but marked as \"not analyzed\". This number does not count toward the daily reserve)", result.SubmissionExemptionCount)

	s.reserve.Set(reserveCount)
	s.submission.Set(submissionCount)
	s.exemption.Set(exemptionCount)
	s.remaining.Set(remainingCount)
}

func (s *QuotaWindow) Show(enableMenuItem func()) {
	s.win.Show()
	go s.Update()
}
