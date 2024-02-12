/*
TunnelEffect (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

main.go

Installer main file
*/
package main

import (
	"fmt"
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
	close, err := logging.NewFileLog(InstallLogFolder(), installWizardLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewFileLog: %v", err)
		os.Exit(10)
	}
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()
	logging.Infof("Start. Version %s Build %s", globals.Version, globals.Build)
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
