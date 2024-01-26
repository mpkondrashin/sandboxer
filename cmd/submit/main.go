package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"examen/pkg/config"
	"examen/pkg/fifo"
	"examen/pkg/globals"
	"examen/pkg/logging"
)

const submitLog = "submit.log"

var (
	fifoMissingDarwinPrefix = "open/create fifo failed"
	fifoMissingDarwinSuffix = "device not configured"
	fifoMisingWindows       = "create file failed: The system cannot find the file specified."
)

func ExamenIsDown(err error) bool {
	if runtime.GOOS == "darwin" {
		return strings.HasPrefix(err.Error(), fifoMissingDarwinPrefix) &&
			strings.HasSuffix(err.Error(), fifoMissingDarwinSuffix)
	}
	if runtime.GOOS == "windows" {
		return strings.HasPrefix(err.Error(), fifoMisingWindows)
	}
	return false
}

func LaunchExamen(conf *config.Configuration) {
	logging.Infof("Launch Examen")
	examenFileName := "examen"
	if runtime.GOOS == "windows" {
		examenFileName += ".exe"
	}
	examenPath := filepath.Join(conf.Folder, examenFileName)
	//examenPath = "../examen/examen"
	cmd := exec.Command(examenPath)
	err := cmd.Start()
	if err != nil {
		//logging.Errorf("%v", err)
		panic(err)
	}
	logging.Infof("Launched Examen")
}

func OpenFIFO(conf *config.Configuration) *fifo.Writer {
	fifoWriter, err := fifo.NewWriter()
	if err == nil {
		return fifoWriter
	}
	if !ExamenIsDown(err) {
		panic(err)
	}
	LaunchExamen(conf)
	for i := 0; i < 10; i++ {
		logging.Debugf("Wait for FIFO")
		fifoWriter, err = fifo.NewWriter()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fifoWriter
}

func main() {
	configFilePath, err := config.FilePath()
	if err != nil {
		panic(err)
	}
	conf := config.New(configFilePath)
	if err := conf.Load(); err != nil {
		panic(err)
	}
	//close := logging.NewFileLog(conf.LogFolder(), submitLog)
	logFolder, err := conf.LogFolder()
	if err != nil {
		panic(err)
	}
	close := logging.NewFileLog(logFolder, submitLog)
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()
	logging.Infof("%s Version %s Submit Started", globals.AppName, globals.Version)
	if len(os.Args) != 2 {
		logging.Errorf("Missing or wrong number of parameters: %s", strings.Join(os.Args[1:], " "))
		return
	}
	filePath := os.Args[1]
	logging.Infof("Submit \"%s\"", filePath)
	fifoWriter := OpenFIFO(conf)
	if fifoWriter == nil {
		panic("did not ")
	}
	defer func() {
		logging.LogError(fifoWriter.Close())
	}()
	if err = fifoWriter.Write(filePath); err != nil {
		logging.Errorf("%v", err)
		log.Println(err)
		os.Exit(2)
	}

	logging.Infof("Submit finished")
}
