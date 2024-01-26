package main

import (
	"log"
	"os"
	"strings"

	"examen/pkg/config"
	"examen/pkg/fifo"
	"examen/pkg/globals"
	"examen/pkg/logging"
)

const submitLog = "submit.log"

var fifoMissing = "open/create fifo failed: open /tmp/examen_fifo: device not configured"

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
	fifoWriter, err := fifo.NewWriter()
	if err != nil {
		logging.Errorf("%v (%T)", err, err)
		log.Println(err)
		os.Exit(1)
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
