package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"sandboxer/pkg/config"
	"sandboxer/pkg/dispatchers"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
)

type Status interface {
	Get(from int, count int)
}

type SandboxerApp struct {
	app                fyne.App
	menu               *fyne.Menu
	submissionMenuItem *fyne.MenuItem
	quotaMenuItem      *fyne.MenuItem
	optionsMenuItem    *fyne.MenuItem
	submissionsWindow  *SubmissionsWindow
	quotaWindow        *QuotaWindow
	optionsWindow      *OptionsWindow
}

func NewSandboxingApp(conf *config.Configuration, channels *dispatchers.Channels, list *task.TaskList) *SandboxerApp {
	fyneApp := app.New()
	deskApp, ok := fyneApp.(desktop.App)
	if !ok {
		panic("only desktop supported")
	}
	a := &SandboxerApp{
		app: fyneApp,
		//submissionsWindow:

	}
	a.quotaWindow = NewQuotaWindow(NewModalWindow(
		a.app.NewWindow("Quota"), a.EnableQuotaMenuItem),
		conf,
	)
	a.submissionsWindow = NewSubmissionsWindow(NewModalWindow(
		a.app.NewWindow("Submissions"), a.EnableSubmissionsMenuItem),
		channels,
		list,
	)
	a.optionsWindow = NewOptionsWindow(
		NewModalWindow(a.app.NewWindow("Options..."), a.EnableOptionsMenuItem),
		conf,
	)
	a.submissionMenuItem = fyne.NewMenuItem("Submissions...", a.Submissions)
	a.quotaMenuItem = fyne.NewMenuItem("Quota...", a.Quota)
	a.optionsMenuItem = fyne.NewMenuItem("Options...", a.Options)
	//	a.submissionMenuItem.Disabled = true
	a.menu = a.Menu()
	deskApp.SetSystemTrayIcon(a.Icon())
	deskApp.SetSystemTrayMenu(a.menu)
	return a
}

func (a *SandboxerApp) Icon() fyne.Resource {
	return ApplicationIcon
}

func (s *SandboxerApp) Run() {
	if len(os.Args) == 2 && os.Args[1] == "--submissions" {
		s.submissionMenuItem.Disabled = true
		s.submissionsWindow.Show()
	}
	s.app.Run()
	//s.quotaWindow.win.ShowAndRun()
}

func (s *SandboxerApp) Menu() *fyne.Menu {
	return fyne.NewMenu(globals.AppName,
		s.submissionMenuItem,
		s.quotaMenuItem,
		s.optionsMenuItem,
		fyne.NewMenuItemSeparator(), //fyne.NewMenuItem("About...", nil),
		fyne.NewMenuItem("Quit", s.Quit),
	)
}

func (s *SandboxerApp) Submissions() {
	s.submissionMenuItem.Disabled = true
	s.menu.Refresh()
	s.submissionsWindow.Show()
}

func (s *SandboxerApp) EnableSubmissionsMenuItem() {
	s.submissionMenuItem.Disabled = false
	s.menu.Refresh()
}

func (s *SandboxerApp) EnableQuotaMenuItem() {
	s.quotaMenuItem.Disabled = false
	s.menu.Refresh()
}

func (s *SandboxerApp) EnableOptionsMenuItem() {
	s.optionsMenuItem.Disabled = false
	s.menu.Refresh()
}

func (s *SandboxerApp) Quota() {
	s.quotaMenuItem.Disabled = true
	s.menu.Refresh()
	s.quotaWindow.Show(s.EnableQuotaMenuItem)
}

func (s *SandboxerApp) Options() {
	s.optionsMenuItem.Disabled = true
	s.menu.Refresh()
	s.optionsWindow.Show()
}

func (s *SandboxerApp) Quit() {
	s.app.Quit()
}

func IconPath(s task.State) string {
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

func HandleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logging.Debugf("Got")
		os.Exit(1)
	}()
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
	//HandleSignals()

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
