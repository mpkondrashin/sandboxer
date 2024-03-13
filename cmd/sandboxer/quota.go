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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
)

type QuotaWindow struct {
	win             fyne.Window
	conf            *config.Configuration
	reserveLabel    *widget.Label //    binding.String
	submissionLabel *widget.Label //    binding.String binding.String
	exemptionLabel  *widget.Label //    binding.String  binding.String
	remainingLabel  *widget.Label //    binding.String  binding.String
	//reserve         binding.String
	//submission      binding.String
	//exemption       binding.String
	//remaining       binding.String
}

func NewQuotaWindow(conf *config.Configuration) *QuotaWindow {
	w := &QuotaWindow{
		conf:            conf,
		reserveLabel:    widget.NewLabel(""), //.NewString(),
		submissionLabel: widget.NewLabel(""), //.NewString(),
		exemptionLabel:  widget.NewLabel(""), //.NewString(),
		remainingLabel:  widget.NewLabel(""), //.NewString(),
	}
	return w
}

func (w *QuotaWindow) Content(modal *ModalWindow) fyne.CanvasObject {
	w.win = modal.win
	reserveCountItem := widget.NewFormItem("Daily Reserve:", w.reserveLabel)                   //widget.NewLabelWithData(w.reserve))
	submissionCountItem := widget.NewFormItem("Files Submitted:", w.submissionLabel)           // widget.NewLabelWithData(w.submission))
	exemptionCountItem := widget.NewFormItem("Unsupported Files Submitted:", w.exemptionLabel) // widget.NewLabelWithData(w.exemption))
	remainingCountItem := widget.NewFormItem("Remaining:", w.remainingLabel)                   // widget.NewLabelWithData(w.remaining))
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
	s.reserveLabel.SetText("?")
	s.submissionLabel.SetText("?")
	s.exemptionLabel.SetText("?")
	s.remainingLabel.SetText("?")
}

func (s *QuotaWindow) Update() {
	logging.Debugf("Run Vision One Quota update")
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

	s.reserveLabel.SetText(reserveCount)
	s.submissionLabel.SetText(submissionCount)
	s.exemptionLabel.SetText(exemptionCount)
	s.remainingLabel.SetText(remainingCount)
}

func (s *QuotaWindow) Show() {
	go s.Update()
}

func (s *QuotaWindow) Hide() {}
