/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

main.go

Sandboxer main file
*/
package main

import (
	"errors"
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
	"sandboxer/pkg/fatal"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
	"sandboxer/pkg/update"
)

type TrayApp struct {
	app  fyne.App
	menu *fyne.Menu
}

type SandboxerApp struct {
	TrayApp
	submissionsWindow *ModalWindow
	optionsWindow     *ModalWindow
	updateWindow      *ModalWindow
}

func NewSandboxingApp(conf *config.Configuration, channels *task.Channels, list *task.TaskList) *SandboxerApp {
	fyneApp := app.New()
	deskApp, ok := fyneApp.(desktop.App)
	if !ok {
		panic("only desktop supported")
	}
	a := &SandboxerApp{
		TrayApp: TrayApp{app: fyneApp},
	}
	submitWindow := NewModalWindow(NewSubmitURLWindow(list, channels), &a.TrayApp)
	quotaWindow := NewModalWindow(NewQuotaWindow(conf), &a.TrayApp)
	//	if conf.SandboxType != config.SandboxVisionOne {
	//		quotaWindow.MenuItem.Disabled = true
	//	}
	statsWindow := NewModalWindow(NewStatsWindow(conf), &a.TrayApp)
	a.submissionsWindow = NewModalWindow(NewSubmissionsWindow(
		channels,
		list,
		conf,
	), &a.TrayApp)

	a.updateWindow = NewModalWindow(NewUpdateWindow(), &a.TrayApp)
	aboutWindow := NewModalWindow(NewAboutWindow(), &a.TrayApp)
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
	a.optionsWindow = NewModalWindow(NewOptionsWindow(conf), &a.TrayApp)

	a.menu = fyne.NewMenu(globals.AppName,
		// SUBMIT_FILE s.submitMenuItem,
		submitWindow.MenuItem,
		a.submissionsWindow.MenuItem,
		fyne.NewMenuItemSeparator(),
		quotaWindow.MenuItem,
		statsWindow.MenuItem,
		a.optionsWindow.MenuItem,
		fyne.NewMenuItemSeparator(),
		a.updateWindow.MenuItem,
		aboutWindow.MenuItem,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", a.Quit),
	)
	deskApp.SetSystemTrayIcon(a.Icon())
	deskApp.SetSystemTrayMenu(a.menu)
	return a
}

func (a *SandboxerApp) Icon() fyne.Resource {
	return ApplicationIcon
}

func (s *SandboxerApp) ShowOptions() {
	confPath, err := globals.ConfigurationFilePath()
	if err != nil {
		logging.Errorf("ConfigurationFilePath: %v", err)
		return
	}
	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		logging.Errorf("Configuration is missing: %v", err)
		s.optionsWindow.Show()
	}
}

func (s *SandboxerApp) Run() {
	if len(os.Args) == 2 && os.Args[1] == "--submissions" {
		s.submissionsWindow.Show()
	}
	s.ShowOptions()
	go s.CheckUpdate()
	s.app.Run()
}

func (s *SandboxerApp) CheckUpdate() {
	logging.Debugf("Run check update")
	need, err := update.NeedUpdateWindow()
	if err != nil {
		logging.LogError(err)
		return
	}
	if need {
		s.updateWindow.Show()
	}
}

func (s *SandboxerApp) Menu() *fyne.Menu {
	return nil
}

func (s *SandboxerApp) SubmitFile() {

}

func (s *SandboxerApp) Quit() {
	s.app.Quit()
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
		logging.Debugf("Run wait for signal")
		<-c
		logging.Debugf("Got signal")
		os.Exit(globals.ExitGotSignal)
	}()
}

func main() {
	configFilePath, err := globals.ConfigurationFilePath()
	if err != nil {
		msg := fmt.Sprintf("globals.ConfigurationFilePath: %v", err)
		fmt.Fprintln(os.Stderr, msg)
		fatal.Warning("Configuration Path Error", msg)
		os.Exit(globals.ExitGetConfigurationFileathError)
	}
	close, err := globals.SetupLogging(globals.Name + ".log")
	if err != nil {
		fmt.Println(err)
	} else {
		defer close()
	}
	conf := config.New(configFilePath)
	if err := conf.Load(); err != nil {
		logging.LogError(err)
		fmt.Fprintf(os.Stderr, "conf.Load: %v", err)
		if !errors.Is(err, os.ErrNotExist) {
			fatal.Warning("Config Error", err.Error())
			os.Exit(globals.ExitLoadConfigError)
		}
	}

	fontFileName := "DroidSansHebrew-Regular.ttf"
	os.Setenv("FYNE_FONT", conf.Resource(fontFileName))
	//	logging.Debugf("FONT: %s", conf.Resource(fontFileName))
	logging.Infof("%s Version %s Build %s Start", globals.AppName, globals.Version, globals.Build)
	logging.Debugf("Configuration file: %s", configFilePath)
	removePid, err := SavePid()
	if err != nil {
		msg := fmt.Sprintf("Save pid: %v", err)
		logging.Errorf(msg)
		fmt.Fprintln(os.Stderr, msg)
		fatal.Warning("Save PID Error", msg)
		os.Exit(globals.ExitSavePidError)
	}
	defer removePid()
	//logging.LogError(ExtractService())
	//list := task.NewList()
	channels := task.NewChannels()
	list := task.NewList()
	launcher := dispatchers.NewLauncher(conf, channels, list)
	launcher.Run()
	defer launcher.Stop()
	app := NewSandboxingApp(conf, channels, list)
	app.Run()
}
