package main

import (
	"os"
	"path/filepath"
	"runtime"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
)

const installWizardLog = globals.Name + "_setup_wizard.log"

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func InstallLogFolder() string {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(path)
}
func main() {
	close := logging.NewFileLog(InstallLogFolder(), installWizardLog)
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
