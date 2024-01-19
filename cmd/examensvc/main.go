package main

import (
	"examen/pkg/config"
	"examen/pkg/globals"
	"examen/pkg/logging"
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"
)

type ExamenSvc struct {
	stop func()
}

// NewExamenSvc - create new sevice
func NewExamenSvc() *ExamenSvc {
	return &ExamenSvc{}
}

// Start - start ExamenSvc service
func (t *ExamenSvc) Start(s service.Service) error {
	logging.Infof("Start Examen")
	//tl.Printf("TunnelEffect Start(%v)", s)
	var err error
	t.stop, err = RunService()
	//tl.Printf("after Run(): stop = %v, err = %v", t.stop, err)
	return err
}

// Stop - stop TunnelEffect service
func (t *ExamenSvc) Stop(s service.Service) error {
	logging.Infof("Stop Examen Service")
	if t.stop == nil {
		return fmt.Errorf("stop is nil")
	}
	t.stop()
	return nil
}

func main() {
	config, err := config.LoadConfiguration(globals.AppID, globals.ConfigFileName)
	if err != nil {
		panic(err)
	}
	close := logging.NewFileLog(config.LogFolder(), examenSvcLog)
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()
	tes := NewExamenSvc()
	s, err := config.Service(tes)
	//	tl.Printf("service.New(): %v, %v", s, err)
	if err != nil {
		log.Println(err)
		//os.Exit(exitcode.ServiceCreate)
		os.Exit(99)
	}
	err = s.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(96)
	}
}
