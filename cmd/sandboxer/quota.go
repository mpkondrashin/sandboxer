package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/mpkondrashin/vone"

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

func NewQuotaWindow(app fyne.App, conf *config.Configuration) *QuotaWindow {
	s := &QuotaWindow{
		win:        app.NewWindow("Quota"),
		conf:       conf,
		reserve:    binding.NewString(),
		submission: binding.NewString(),
		exemption:  binding.NewString(),
		remaining:  binding.NewString(),
	}
	s.Reset()
	s.win.SetCloseIntercept(func() {
		//logging.Debugf("XXX Close")
		s.Hide()
	})
	//	s.win.Resize(fyne.Size{Width: 400, Height: 300})
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
	//s.win.Content().MinSize (fyne.Size{Width: 400, Height: 300})
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
	return
	reserveCountItem := widget.NewFormItem("Daily Reserve:", widget.NewLabel(reserveCount))
	submissionCountItem := widget.NewFormItem("Files Submitted:", widget.NewLabel(submissionCount))
	exemptionCountItem := widget.NewFormItem("Unsupported Files Submitted:", widget.NewLabel(exemptionCount))
	remainingCountItem := widget.NewFormItem("Remaining:", widget.NewLabel(remainingCount))
	form := widget.NewForm(
		reserveCountItem,
		submissionCountItem,
		exemptionCountItem,
		remainingCountItem,
	)
	s.win.SetContent(form)
	//str := binding.NewString()
	/*
		SubmissionReserveCount   int `json:"submissionReserveCount"`
		SubmissionRemainingCount int `json:"submissionRemainingCount"`
		SubmissionCount          int `json:"submissionCount"`
		SubmissionExemptionCount int `json:"submissionExemptionCount"`
		SubmissionCountDetail    struct {
			FileCount          int `json:"fileCount"`
			FileExemptionCount int `json:"fileExemptionCount"`
			URLCount           int `json:"urlCount"`
			URLExemptionCount  int `json:"urlExemptionCount"`
		} `json:"submissionCountDetail"`
	*/
}

func (s *QuotaWindow) Show() {
	s.win.Show()
	go s.Update()
}

func (s *QuotaWindow) Hide() {
	s.win.Hide()
}
