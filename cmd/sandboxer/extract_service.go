package main

import (
	"embed"
	"os"
	"path/filepath"

	"sandboxer/pkg/extract"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
)

//go:embed embed/*
var embedFS embed.FS

var (
	ServicesFolder = "Library/Services"
	WorkflowFolder = globals.AppName + ".workflow"
	//ServicePath    = ServicesFolder + "/" + WorkflowFolder
	ServiceTGZPath = "embed/" + globals.AppName + ".tar.gz"
)

func serviceExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ExtractService() error {
	homeFolder, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	servicesPath := filepath.Join(homeFolder, ServicesFolder)
	if serviceExist(filepath.Join(servicesPath, WorkflowFolder)) {
		logging.Infof("Extract Service: Service Exist")
		return nil
	}
	logging.Infof("Extract Service")
	return extract.ExtractFileTGZ(servicesPath, embedFS, ServiceTGZPath)
}
