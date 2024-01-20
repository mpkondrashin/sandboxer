package main

import (
	"examen/pkg/config"
	"examen/pkg/globals"
	"examen/pkg/logging"
	"log"
	"os"

	"github.com/kardianos/service"
)

var logger service.Logger

type ExamenSvc struct {
	stop func()
}

// NewExamenSvc - create new sevice
func NewExamenSvc() *ExamenSvc {
	return &ExamenSvc{}
}

// Start - start ExamenSvc service
func (t *ExamenSvc) Start(s service.Service) error {
	logging.Infof("Start Examen") // XXXX
	logger.Infof("Start %s", globals.AppName)
	var err error
	t.stop, err = RunService()
	return err
}

// Stop - stop TunnelEffect service
func (t *ExamenSvc) Stop(s service.Service) error {
	logger.Info("Stop %s Service", globals.AppName)
	logging.Infof("Stop Examen Service") // XXXX
	if t.stop == nil {
		return logger.Info("stop is nil")
	}
	t.stop()
	return nil
}

func main() {
	svc, err := service.New(nil, &service.Config{Name: globals.SvcName})
	if err != nil {
		log.Fatal(err)
	}
	logger, err = svc.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	config, err := config.LoadConfiguration(globals.AppID, globals.ConfigFileName)
	if err != nil {
		logger.Errorf("Configuration Load Error: %v", err)
	}

	close := logging.NewFileLog(config.LogFolder(), examenSvcLog)
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("panic Error: %v", err)
			logging.Criticalf("panic: %v", err)
		}
	}()

	tes := NewExamenSvc()
	s, err := config.Service(tes)
	if err != nil {
		logging.Errorf("Create Service Error: %v", err)
		logging.Criticalf("Create Service Error: %v", err)
		os.Exit(99)
	}

	err = s.Run()
	logging.Debugf("Run() returned: %v", err)
	if err != nil {
		logging.Errorf("Run Service Error: %v", err)
		logging.Criticalf("Run Service Error: %v", err)
		os.Exit(96)
	}
}
