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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/ddan"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
)

type StatStruct struct {
	Last4Hours  binding.String
	Last24Hours binding.String
	Last7Days   binding.String
	Last30Days  binding.String
	Last90Days  binding.String
}

func NewStatStruct() *StatStruct {
	return &StatStruct{
		Last4Hours:  binding.NewString(),
		Last24Hours: binding.NewString(),
		Last7Days:   binding.NewString(),
		Last30Days:  binding.NewString(),
		Last90Days:  binding.NewString(),
	}
}

func (s *StatStruct) FormItems() []*widget.FormItem {
	return []*widget.FormItem{
		widget.NewFormItem("Last 4 Hours:", widget.NewLabelWithData(s.Last4Hours)),
		widget.NewFormItem("Last 24 Hours:", widget.NewLabelWithData(s.Last24Hours)),
		widget.NewFormItem("Last 7 Days:", widget.NewLabelWithData(s.Last7Days)),
		widget.NewFormItem("Last 30 Days:", widget.NewLabelWithData(s.Last30Days)),
		widget.NewFormItem("Last 90 Days:", widget.NewLabelWithData(s.Last90Days)),
	}
}

func (s *StatStruct) Reset() {
	s.Last4Hours.Set("?")
	s.Last24Hours.Set("?")
	s.Last7Days.Set("?")
	s.Last30Days.Set("?")
	s.Last90Days.Set("?")
}

func toString(v int) string {
	if v == -1 {
		return "no data"
	}
	return fmt.Sprintf("%d sec", v)
}

func (s *StatStruct) Set(t *ddan.AvgTime) {
	s.Last4Hours.Set(toString(t.Last4Hours))
	s.Last24Hours.Set(toString(t.Last24Hours))
	s.Last7Days.Set(toString(t.Last7Days))
	s.Last30Days.Set(toString(t.Last30Days))
	s.Last90Days.Set(toString(t.Last90Days))
}

type StatsWindow struct {
	win  fyne.Window
	conf *config.Configuration

	AvgVAAnalysisTime      *StatStruct
	AvgTotalProcessingTime *StatStruct
}

func NewStatsWindow(conf *config.Configuration) *StatsWindow {
	w := &StatsWindow{
		conf:                   conf,
		AvgVAAnalysisTime:      NewStatStruct(),
		AvgTotalProcessingTime: NewStatStruct(),
	}
	return w
}

func (w *StatsWindow) Content(modal *ModalWindow) fyne.CanvasObject {
	w.win = modal.win
	totalLabel := widget.NewLabel("Average Total Processing Time\n(from the moment of submission,\nincluding queueing)")
	totalForm := widget.NewForm(w.AvgTotalProcessingTime.FormItems()...)
	vaLabel := widget.NewLabel("Average Analysis Time\n(by sandbox)\n")
	vaForm := widget.NewForm(w.AvgVAAnalysisTime.FormItems()...)
	totalVBox := container.NewVBox(totalLabel, totalForm)
	vaVBox := container.NewVBox(vaLabel, vaForm)
	dataHBox := container.NewHBox(totalVBox, vaVBox)
	w.Reset()
	return container.NewVBox(
		dataHBox,
		widget.NewButton("Update", func() {
			w.Update()
		}),
	)
}

func (w *StatsWindow) Name() string {
	return "Analyzer Statistics"
}

func (s *StatsWindow) Reset() {
	s.AvgTotalProcessingTime.Reset()
	s.AvgVAAnalysisTime.Reset()
}

func (s *StatsWindow) Update() {
	s.Reset()
	analyzer, err := s.conf.DDAn.AnalyzerWithUUID()
	if err != nil {
		logging.LogError(err)
		dialog.ShowError(err, s.win)
		return
	}
	stats, err := analyzer.GetStats(context.TODO())
	if err != nil {
		logging.LogError(err)
		dialog.ShowError(err, s.win)
		return
	}
	s.AvgTotalProcessingTime.Set(&stats.AvgTotalProcessingTime)
	s.AvgVAAnalysisTime.Set(&stats.AvgVAAnalysisTime)
}

func (s *StatsWindow) Show() {
	go s.Update()
}

func (s *StatsWindow) Hide() {}
