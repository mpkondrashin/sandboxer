/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

main.go

Send file to sandboxer
*/
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"sandboxer/pkg/config"
	"sandboxer/pkg/fatal"
	"sandboxer/pkg/fifo"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/xplatform"
)

const submitLog = "submit.log"

var ErrUnsupportedOS = errors.New("unsupported OS")

func SubmissionsExecutablePath(conf *config.Configuration) (string, error) {
	return xplatform.ExecutablePath(conf.GetFolder(), globals.AppName, globals.Name)
}

func LaunchSandboxer(conf *config.Configuration) {
	logging.Infof("Launch " + globals.AppName)
	executablePath, err := SubmissionsExecutablePath(conf)
	if err != nil {
		fatal.Warning("Config Error", err.Error())
		panic(err)
	}
	logging.Infof("Run " + executablePath)
	cmd := exec.Command(executablePath, "--submissions")
	if err := cmd.Start(); err != nil {
		fatal.Warning("Execute Error", err.Error())
		panic(err)
	}
	//	if err := cmd.Wait(); err != nil {
	//		logging.Debugf(globals.AppName+" exited with with error: %v", err)
	//	}
	logging.Infof("Launched " + globals.AppName)
}

func OpenFIFO(conf *config.Configuration) *fifo.Writer {
	fifoWriter, err := fifo.NewWriter()
	if err == nil {
		return fifoWriter
	}
	if !fifo.IsDown(err) {
		fatal.Warning("FIFO Error", err.Error())
		panic(err)
	}
	LaunchSandboxer(conf)
	for i := 0; i < 10; i++ {
		logging.Debugf("Wait for FIFO")
		fifoWriter, err = fifo.NewWriter()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return fifoWriter
}

func main() {
	configFilePath, err := globals.ConfigurationFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ConfigurationFilePath: %v", err)
		fatal.Warning("Configuration Error", err.Error())
		os.Exit(10)
	}
	conf := config.New(configFilePath)
	if err := conf.Load(); err != nil {
		if runtime.GOOS == "windows" { // After creating darwin installer this if should be removed
			fmt.Fprintf(os.Stderr, "conf.Load: %v", err)
			fatal.Warning("Configuration Error", err.Error())
			os.Exit(20)
		}
	}
	//close := logging.NewFileLog(conf.LogFolder(), submitLog)
	closeLogging, err := globals.SetupLogging(submitLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SetupLogging: %v", err)
		fatal.Warning("Setup Logging Error", err.Error())
		os.Exit(30)
	}
	defer closeLogging()
	defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()
	logging.Infof("%s Version %s Build %s Submit Started", globals.AppName, globals.Version, globals.Build)
	if len(os.Args) != 2 {
		msg := fmt.Sprintf("Missing or wrong number of parameters: %s", strings.Join(os.Args[1:], " "))
		logging.Errorf(msg)
		fatal.Warning("Arguments Error", msg)
		os.Exit(40)
	}
	filePath := os.Args[1]
	logging.Infof("Submit \"%s\"", filePath)
	fifoWriter := OpenFIFO(conf)
	if fifoWriter == nil {
		msg := fmt.Sprintf(globals.AppName + " is not running and can not be launched")
		logging.Errorf(msg)
		fatal.Warning("Configuration Error", msg)
		os.Exit(50)
	}
	defer func() {
		logging.LogError(fifoWriter.Close())
	}()
	if err = fifoWriter.Write(filePath); err != nil {
		msg := fmt.Sprintf("fifoWriter.Write: %v", err)
		logging.Errorf(msg)
		fatal.Warning("Configuration Error", msg)
		os.Exit(60)
	}
	logging.Infof("Submit finished")
}
