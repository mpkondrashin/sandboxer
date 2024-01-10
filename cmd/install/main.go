package main

import (
	"os"
	"runtime"

	"examen/pkg/logging"
)

const (
	appName          = "Examen"
	configFileName   = "examen.yaml" // remove - use fyne
	appID            = "com.github.mpkondrashin.examen"
	installWizardLog = "examen_setup_wizard.log"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func main() {
	close := logging.NewFileLog(logging.InstallLogFolder(), installWizardLog)
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()
	logging.Infof("Start")
	logging.Debugf("OS: %s (%s)", runtime.GOOS, runtime.GOARCH)
	capturesFolder := ""
	if len(os.Args) == 3 && os.Args[1] == "--capture" {
		capturesFolder = os.Args[2]
	}
	logging.Infof("Starting Wizard")
	c := NewWizard(capturesFolder)
	c.Run()
	logging.Infof("Setup finished")
}
