package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"examen/pkg/config"
	"examen/pkg/logging"
	"examen/pkg/state"
	"examen/pkg/submit"
	"examen/pkg/task"
)

type Status interface {
	Get(from int, count int)
}

type ExamenApp struct {
	app               fyne.App
	submissionsWindow *SubmissionsWindow
	optionsWindow     *OptionsWindow
}

func NewExamenApp(conf *config.Configuration) *ExamenApp {
	fyneApp := app.New()
	deskApp, ok := fyneApp.(desktop.App)
	if !ok {
		panic("only desktop supported")
	}
	a := &ExamenApp{
		app:               fyneApp,
		submissionsWindow: NewSubmissionsWindow(fyneApp),
		optionsWindow:     NewOptionsWindow(fyneApp, conf),
	}
	deskApp.SetSystemTrayIcon(a.Icon())
	deskApp.SetSystemTrayMenu(a.Menu())
	return a
}

func (a *ExamenApp) Icon() fyne.Resource {
	path := "../../resources/LowRisk.svg"
	r, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		panic(err)
	}
	return r
}

func (s *ExamenApp) Run() {
	s.submissionsWindow.Show()
	s.submissionsWindow.win.ShowAndRun()
}

func (s *ExamenApp) Menu() *fyne.Menu {
	return fyne.NewMenu("Examen",
		fyne.NewMenuItem("Submissions...", s.Submissions),
		fyne.NewMenuItem("Options...", s.Options),
		//fyne.NewMenuItem("About...", nil),
		fyne.NewMenuItem("Quit", s.Quit),
	)
}

func (s *ExamenApp) Submissions() {
	s.submissionsWindow.Show()
}

func (s *ExamenApp) Options() {
	s.optionsWindow.Show()
}

func (s *ExamenApp) Quit() {
	s.app.Quit()
}

func IconPath(s state.State) string {
	return fmt.Sprintf("../../resources/%s.svg", s.String())
}

var Conf *config.Configuration

func main() {
	configFilePath, err := config.FilePath()
	if err != nil {
		panic(err)
	}
	conf := config.New(configFilePath)

	close := setupLogging()
	defer close()

	list := task.NewList()
	stop, err := submit.RunService(conf, list)
	if err != nil {
		panic(err)
	}
	defer stop()
	app := NewExamenApp(conf)
	app.Run()
}

func setupLogging() func() {
	logging.SetLevel(logging.DEBUG)
	//      logFileName := fmt.Sprintf("setup_%s.log", time.Now().Format("20060102_150405"))
	logFileName := "examen.log"
	logFilePath := filepath.Join(".", logFileName)
	//      logFilePath := filepath.Join(os.TempDir(), logFileName)
	fmt.Println(logFilePath)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logger := logging.NewFileLogger(file)
	logging.AddLogger(logger)
	return func() {
		logging.Infof("Close Logging")
		logging.Close()
		file.Close()
	}
}
