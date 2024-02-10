package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"sandboxer/pkg/config"
	"sandboxer/pkg/fifo"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
)

const submitLog = "submit.log"

var (
	fifoMissingDarwinPrefix = "open/create fifo failed"
	fifoMissingDarwinSuffix = "device not configured"
	fifoMisingWindows       = "create file failed: The system cannot find the file specified."
)

func IsDown(err error) bool {
	if runtime.GOOS == "darwin" {
		return strings.HasPrefix(err.Error(), fifoMissingDarwinPrefix) &&
			strings.HasSuffix(err.Error(), fifoMissingDarwinSuffix)
	}
	if runtime.GOOS == "windows" {
		return strings.HasPrefix(err.Error(), fifoMisingWindows)
	}
	return false
}

var ErrUnsupportedOS = errors.New("unsupported OS")

func SubmissionsExecutablePath(conf *config.Configuration) (string, error) {
	if runtime.GOOS == "windows" {
		return filepath.Join(conf.Folder, globals.AppName, globals.AppName+".exe"), nil
	}
	if runtime.GOOS == "darwin" {
		return fmt.Sprintf("%s/%s.app/Contents/MacOS/%s", conf.Folder, globals.AppName, globals.Name), nil
	}
	return "", fmt.Errorf("%s: %W", runtime.GOOS, ErrUnsupportedOS)
}

func LaunchSandboxer(conf *config.Configuration) {
	logging.Infof("Launch " + globals.AppName)
	executableFileName := globals.AppName
	if runtime.GOOS == "windows" {
		executableFileName += ".exe"
	}
	executablePath, err := SubmissionsExecutablePath(conf) // := filepath.Join(conf.Folder, globals.AppName, executableFileName)
	if err != nil {
		panic(err)
	}
	logging.Infof("Run " + executablePath)
	cmd := exec.Command(executablePath, "--submissions")
	if err := cmd.Start(); err != nil {
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
	if !IsDown(err) {
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
		os.Exit(10)
	}
	conf := config.New(configFilePath)
	if err := conf.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "conf.Load: %v", err)
		os.Exit(20)
	}
	//close := logging.NewFileLog(conf.LogFolder(), submitLog)
	closeLogging, err := globals.SetupLogging(submitLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SetupLogging: %v", err)
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
		logging.Errorf("Missing or wrong number of parameters: %s", strings.Join(os.Args[1:], " "))
		os.Exit(40)
	}
	filePath := os.Args[1]
	logging.Infof("Submit \"%s\"", filePath)
	fifoWriter := OpenFIFO(conf)
	if fifoWriter == nil {
		logging.Errorf(globals.AppName + " is not running and can not be launched")
		os.Exit(50)
	}
	defer func() {
		logging.LogError(fifoWriter.Close())
	}()
	if err = fifoWriter.Write(filePath); err != nil {
		logging.Errorf("fifoWriter.Write: %v", err)
		os.Exit(60)
	}
	logging.Infof("Submit finished")
}
