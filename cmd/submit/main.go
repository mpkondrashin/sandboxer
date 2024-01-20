package main

import (
	"os"
	"strings"

	"examen/pkg/config"
	"examen/pkg/fifo"
	"examen/pkg/globals"
	"examen/pkg/logging"
)

const submitLog = "submit.log"

func main() {
	config, err := config.LoadConfiguration(globals.AppID, globals.ConfigFileName)
	if err != nil {
		panic(err)
	}
	close := logging.NewFileLog(config.LogFolder(), submitLog)
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()
	logging.Infof("Submit")
	if len(os.Args) != 2 {
		logging.Errorf("Wrong parameters: %s", strings.Join(os.Args, " "))
		return
	}
	filePath := os.Args[1]
	logging.Infof("Submit \"%s\"", filePath)
	fifoWriter, err := fifo.NewWriter()
	logging.LogError(err)
	if err != nil {
		os.Exit(1)
	}
	fifoWriter.Write(os.Args[1])

	logging.Infof("Submit finished")
}
