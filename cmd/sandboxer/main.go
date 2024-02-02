package main

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"sandboxer/pkg/config"
	"sandboxer/pkg/dispatchers"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/state"
	"sandboxer/pkg/task"
)

type Status interface {
	Get(from int, count int)
}

type SandboxerApp struct {
	app               fyne.App
	submissionsWindow *SubmissionsWindow
	quotaWindow       *QuotaWindow
	optionsWindow     *OptionsWindow
}

func NewSandboxingApp(conf *config.Configuration, channels *dispatchers.Channels, list *task.TaskList) *SandboxerApp {
	fyneApp := app.New()
	deskApp, ok := fyneApp.(desktop.App)
	if !ok {
		panic("only desktop supported")
	}
	a := &SandboxerApp{
		app:               fyneApp,
		submissionsWindow: NewSubmissionsWindow(fyneApp, channels, list),
		quotaWindow:       NewQuotaWindow(fyneApp, conf),
		optionsWindow:     NewOptionsWindow(fyneApp, conf),
	}
	deskApp.SetSystemTrayIcon(a.Icon())
	deskApp.SetSystemTrayMenu(a.Menu())
	return a
}

func (a *SandboxerApp) Icon() fyne.Resource {
	return ApplicationIcon
	path := "../../resources/LowRisk.svg"
	r, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		panic(err)
	}
	return r
}

func (s *SandboxerApp) Run() {
	s.submissionsWindow.Show()
	s.submissionsWindow.win.ShowAndRun()
}

func (s *SandboxerApp) Menu() *fyne.Menu {
	return fyne.NewMenu(globals.AppName,
		fyne.NewMenuItem("Submissions...", s.Submissions),
		fyne.NewMenuItem("Quota...", s.Quota),
		fyne.NewMenuItem("Options...", s.Options),
		//fyne.NewMenuItem("About...", nil),
		fyne.NewMenuItem("Quit", s.Quit),
	)
}

func (s *SandboxerApp) Submissions() {
	s.submissionsWindow.Show()
}

func (s *SandboxerApp) Quota() {
	s.quotaWindow.Show()
}

func (s *SandboxerApp) Options() {
	s.optionsWindow.Show()
}

func (s *SandboxerApp) Quit() {
	s.app.Quit()
}

func IconPath(s state.State) string {
	return fmt.Sprintf("../../resources/%s.svg", s.String())
}

//var Conf *config.Configuration

func SavePid() (func(), error) {
	pidFilePath, err := globals.PidFilePath()
	if err != nil {
		return nil, err
	}
	pid := strconv.Itoa(os.Getpid())
	if err := os.WriteFile(pidFilePath, []byte(pid), 0644); err != nil {
		return nil, err
	}
	return func() {
		os.Remove(pidFilePath)
	}, nil
}

func main() {
	configFilePath, err := globals.ConfigurationFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "globals.ConfigurationFilePath: %v", err)
		os.Exit(10)
	}
	conf := config.New(configFilePath)
	if err := conf.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "conf.Load: %v", err)
		os.Exit(20)
	}
	close, err := globals.SetupLogging(globals.Name + ".log")
	if err != nil {
		fmt.Println(err)
	} else {
		defer close()
	}
	logging.Infof("%s Version %s Start", globals.AppName, globals.Version)
	logging.Debugf("Configuration file: %s", configFilePath)
	removePid, err := SavePid()
	if err != nil {
		logging.Errorf("Save pid: %v", err)
		fmt.Fprintf(os.Stderr, "Save pid: %v", err)
		os.Exit(30)

	}
	defer removePid()
	//list := task.NewList()
	channels := dispatchers.NewChannels()
	list := task.NewList()
	launcher := dispatchers.NewLauncher(conf, channels, list)
	launcher.Run() //stop, err := submit.RunService(conf, list)
	//if err != nil {
	//	logging.Errorf("RunService: %v", err)
	//	fmt.Fprintf(os.Stderr, "RunService: %v", err)
	//	os.Exit(40)
	//}
	defer launcher.Stop()
	app := NewSandboxingApp(conf, channels, list)
	app.Run()
}
