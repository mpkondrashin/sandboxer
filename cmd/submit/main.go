package main

import (
	"os"
	"strings"

	"examen/pkg/config"
	"examen/pkg/fifo"
	"examen/pkg/logging"
)

const submitLog = "submit.log"

func main() {
	configFilePath, err := config.FilePath()
	if err != nil {
		panic(err)
	}
	conf := config.New(configFilePath)
	_ = conf
	//close := logging.NewFileLog(conf.LogFolder(), submitLog)
	close := logging.NewFileLog(".", submitLog)
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()
	logging.Infof("Submit Started")
	if len(os.Args) != 2 {
		logging.Errorf("Missing or wrong number of parameters: %s", strings.Join(os.Args[1:], " "))
		return
	}
	filePath := os.Args[1]
	logging.Infof("Submit \"%s\"", filePath)
	fifoWriter, err := fifo.NewWriter()
	logging.LogError(err)
	if err != nil {
		os.Exit(1)
	}
	defer fifoWriter.Close()
	err = fifoWriter.Write(os.Args[1])
	logging.LogError(err)
	logging.Infof("Submit finished")
}
