/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
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

	"sandboxer/pkg/fatal"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"

	"github.com/virtuald/go-paniclog"
)

const installWizardLog = globals.Name + "_setup_wizard.log"

// TODO:
// show fatal.Warning in case of error
// change 10 to globals constant

func SetupLogging(logFileName string) (func(), error) {
	logging.SetLevel(logging.DEBUG)
	path, err := os.Executable()
	if err != nil {
		return nil, err
	}
	logFolder := filepath.Dir(path)
	logging.SetLevel(logging.DEBUG)
	file, err := logging.OpenRotated(logFolder, logFileName, 0644, globals.MaxLogFileSize, globals.LogsKeep)
	if err != nil {
		return nil, err
	}
	paniclog.RedirectStderr(file.File)
	logging.SetLogger(logging.NewFileLogger(file))
	return func() {
		logging.Infof("Close Logging")
		file.Close()
	}, nil
}

func main() {
	close, err := SetupLogging(installWizardLog)
	if err != nil {
		msg := fmt.Sprintf("NewFileLog: %v", err)
		fmt.Fprintln(os.Stderr, msg)
		fatal.Warning("SetupLogging Error", msg)
		os.Exit(10)
	}
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	/*defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()*/
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
