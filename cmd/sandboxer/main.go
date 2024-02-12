/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

main.go

Sandboxer main file
*/
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
	app  fyne.App
	menu *fyne.Menu
	// SUBMIT_FILE submitMenuItem     *fyne.MenuItem
	submissionsMenuItem *fyne.MenuItem
	quotaMenuItem       *fyne.MenuItem
	optionsMenuItem     *fyne.MenuItem
	aboutMenuItem       *fyne.MenuItem

	submissionsWindow *SubmissionsWindow
	quotaWindow       *QuotaWindow
	optionsWindow     *OptionsWindow
	aboutWindow       *AboutWindow
}

func NewSandboxingApp(conf *config.Configuration, channels *dispatchers.Channels, list *task.TaskList) *SandboxerApp {
	fyneApp := app.New()
	deskApp, ok := fyneApp.(desktop.App)
	if !ok {
		panic("only desktop supported")
	}
	a := &SandboxerApp{
		app: fyneApp,
	}
	a.quotaWindow = NewQuotaWindow(NewModalWindow(
		a.app.NewWindow("Quota"), a.EnableQuotaMenuItem),
		conf,
	)
	a.submissionsWindow = NewSubmissionsWindow(NewModalWindow(
		a.app.NewWindow("Submissions"), a.EnableSubmissionsMenuItem),
		channels,
		list,
		conf,
	)
	a.optionsWindow = NewOptionsWindow(
		NewModalWindow(a.app.NewWindow("Options"), a.EnableOptionsMenuItem),
		conf,
	)
	a.aboutWindow = NewAboutWindow(
		NewModalWindow(a.app.NewWindow("About"), a.EnableAboutMenuItem),
	)
	/* SUBMIT_FILE
	a.submitMenuItem = fyne.NewMenuItem("Submit File", func() {
		fmt.Println("Submit file")
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			fmt.Println("Open", err)
			if err != nil {
				return
			}
			channels.TaskChannel[dispatchers.ChPrefilter] <- list.NewTask(uri.URI().String())
			list.Updated()
		}, a.optionsWindow.win)
	})
	*/
	a.submissionsMenuItem = fyne.NewMenuItem("Submissions...", a.Submissions)
	a.quotaMenuItem = fyne.NewMenuItem("Quota...", a.Quota)
	a.optionsMenuItem = fyne.NewMenuItem("Options...", a.Options)
	a.aboutMenuItem = fyne.NewMenuItem("About...", a.About)

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
		s.submissionsMenuItem.Disabled = true
		s.submissionsWindow.Show()
	}
	s.app.Run()
}

func (s *SandboxerApp) Menu() *fyne.Menu {
	return fyne.NewMenu(globals.AppName,
		// SUBMIT_FILE s.submitMenuItem,
		s.submissionsMenuItem,
		s.quotaMenuItem,
		s.optionsMenuItem,
		s.aboutMenuItem,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", s.Quit),
	)
}

func (s *SandboxerApp) SubmitFile() {

}

func (s *SandboxerApp) Submissions() {
	s.submissionsMenuItem.Disabled = true
	s.menu.Refresh()
	s.submissionsWindow.Show()
}

func (s *SandboxerApp) EnableSubmissionsMenuItem() {
	s.submissionsMenuItem.Disabled = false
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
func (s *SandboxerApp) EnableAboutMenuItem() {
	s.aboutMenuItem.Disabled = false
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

func (s *SandboxerApp) About() {
	s.aboutMenuItem.Disabled = true
	s.menu.Refresh()
	s.aboutWindow.Show()
}

func (s *SandboxerApp) Quit() {
	s.app.Quit()
}

func IconPath(s task.State) string {
	return fmt.Sprintf("../../resources/%s.svg", s.String())
}

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
		logging.Debugf("Got signal")
		os.Exit(globals.ExitGotSignal)
	}()
}

func main() {
	configFilePath, err := globals.ConfigurationFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "globals.ConfigurationFilePath: %v", err)
		os.Exit(globals.ExitGetConfigurationFileathError)
	}
	conf := config.New(configFilePath)
	if err := conf.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "conf.Load: %v", err)
		os.Exit(globals.ExitLoadConfigError)
	}
	close, err := globals.SetupLogging(globals.Name + ".log")
	if err != nil {
		fmt.Println(err)
	} else {
		defer close()
	}

	logging.Infof("%s Version %s Build %s Start", globals.AppName, globals.Version, globals.Build)
	logging.Debugf("Configuration file: %s", configFilePath)
	removePid, err := SavePid()
	if err != nil {
		logging.Errorf("Save pid: %v", err)
		fmt.Fprintf(os.Stderr, "Save pid: %v", err)
		os.Exit(globals.ExitSavePidError)
	}
	defer removePid()
	//list := task.NewList()
	channels := dispatchers.NewChannels()
	list := task.NewList()
	launcher := dispatchers.NewLauncher(conf, channels, list)
	launcher.Run()
	defer launcher.Stop()
	app := NewSandboxingApp(conf, channels, list)
	app.Run()
}
