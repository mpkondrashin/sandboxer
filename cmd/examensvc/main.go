package main

import (
	"errors"
	"examen/pkg/config"
	"examen/pkg/globals"
	"examen/pkg/logging"
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"
)

type ExamenSvc struct {
	stop    func()
	service service.Service
}

// NewExamenSvc - create new sevice
func NewExamenSvc() *ExamenSvc {
	return &ExamenSvc{}
}

// Start - start TunnelEffect service
func (t *ExamenSvc) Start(s service.Service) error {
	logging.Infof("Start Examen")
	//tl.Printf("TunnelEffect Start(%v)", s)
	var err error
	err = RunService()
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
	//if service.Interactive() {
	//	logger.Info("Exit")
	//	os.Exit(0)
	//}
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

	svcConfig := &service.Config{
		Name:        "examen",
		DisplayName: "Examen",
		Description: "Submit files to Vision One sandbox",
	}
	tes := NewExamenSvc()
	s, err := service.New(tes, svcConfig)
	//	tl.Printf("service.New(): %v, %v", s, err)
	if err != nil {
		log.Fatal(err)
		//os.Exit(exitcode.ServiceCreate)
		os.Exit(99)
	}
	//tes.service = s
	//logger, err = s.Logger(nil)
	//tl.Printf("s.Logger(): %v, %v", logger, err)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 1 {
		err = s.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(96)
		}
		return
	}
	if len(os.Args) == 2 {
		operation := os.Args[1]
		if operation == "status" {
			status, err := s.Status()
			if errors.Is(err, service.ErrNotInstalled) {
				fmt.Println("unknown")
				return
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(87)
			}
			statusName := []string{
				"unknown", "running", "stopped",
			}
			fmt.Println(statusName[status])
			return
		}
		err := service.Control(s, operation)
		if err != nil {
			fmt.Println(err)
			os.Exit(95)
		}
		return
	}
	fmt.Printf("%s {status|install|start|restart|stop|uninstall}", os.Args[0])
	os.Exit(10)
}
