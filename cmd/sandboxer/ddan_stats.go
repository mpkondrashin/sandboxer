/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

ddan_stats.go

Statistics from DDAn
*/
package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/ddan"

	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
)

type StatStruct struct {
	Last4Hours  *widget.Label
	Last24Hours *widget.Label
	Last7Days   *widget.Label
	Last30Days  *widget.Label
	Last90Days  *widget.Label
}

func NewStatStruct() *StatStruct {
	return &StatStruct{
		Last4Hours:  widget.NewLabel(""),
		Last24Hours: widget.NewLabel(""),
		Last7Days:   widget.NewLabel(""),
		Last30Days:  widget.NewLabel(""),
		Last90Days:  widget.NewLabel(""),
	}
}

func (s *StatStruct) FormItems() []*widget.FormItem {
	return []*widget.FormItem{
		widget.NewFormItem("Last 4 Hours:", s.Last4Hours),
		widget.NewFormItem("Last 24 Hours:", s.Last24Hours),
		widget.NewFormItem("Last 7 Days:", s.Last7Days),
		widget.NewFormItem("Last 30 Days:", s.Last30Days),
		widget.NewFormItem("Last 90 Days:", s.Last90Days),
	}
}

func (s *StatStruct) Reset() {
	s.Last4Hours.SetText("?")
	s.Last24Hours.SetText("?")
	s.Last7Days.SetText("?")
	s.Last30Days.SetText("?")
	s.Last90Days.SetText("?")
}

func toString(v int) string {
	if v == -1 {
		return "no data"
	}
	return fmt.Sprintf("%d sec", v)
}

func (s *StatStruct) Set(t *ddan.AvgTime) {
	s.Last4Hours.SetText(toString(t.Last4Hours))
	s.Last24Hours.SetText(toString(t.Last24Hours))
	s.Last7Days.SetText(toString(t.Last7Days))
	s.Last30Days.SetText(toString(t.Last30Days))
	s.Last90Days.SetText(toString(t.Last90Days))
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
	topLabel := widget.NewLabel("Deep Discovery Analyzer Statistis")
	totalLabel := widget.NewLabel("Average Total Processing Time\n(from the moment of submission,\nincluding queueing)")
	totalForm := widget.NewForm(w.AvgTotalProcessingTime.FormItems()...)
	vaLabel := widget.NewLabel("Average Analysis Time\n(by sandbox)\n")
	vaForm := widget.NewForm(w.AvgVAAnalysisTime.FormItems()...)
	totalVBox := container.NewVBox(totalLabel, totalForm)
	vaVBox := container.NewVBox(vaLabel, vaForm)
	dataHBox := container.NewHBox(totalVBox, vaVBox)
	w.Reset()
	return container.NewVBox(
		topLabel,
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
	logging.Debugf("Run DDAn Stats update")
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
