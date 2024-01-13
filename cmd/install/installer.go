package main

import (
	"fmt"
	"os"
	"path/filepath"

	"examen/pkg/config"
	"examen/pkg/globals"
	"examen/pkg/logging"
	"examen/pkg/script"
)

type Installer struct {
	appID string
	//fileName string
	config          config.Configuration
	uninstallScript *os.File
	//hash     string
}

func NewInstaller(appID string) *Installer {
	folder, _ := config.InstallFolder()
	return &Installer{appID: appID,
		config: config.Configuration{Folder: folder}}
}

/*
	func (m *Model) ConfigExists() (bool, error) {
		path, err := m.configFilePath()
		if err != nil {
			return false, err
		}
		_, err = os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
*/
func (i *Installer) LoadConfig() error {
	filePath, err := i.ConfigFileFolder()
	if err != nil {
		return err
	}
	err = i.config.Load(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (i *Installer) Load() error {
	logging.Debugf("Load config")
	folder, err := i.ConfigFileFolder()
	if err != nil {
		return err
	}
	filePath := filepath.Join(folder, globals.ConfigFileName)
	return i.config.Load(filePath)
}

func (i *Installer) Save() error {
	logging.Debugf("Save config")
	folder, err := i.ConfigFileFolder()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(folder, 0700); err != nil {
		return err
	}
	filePath := filepath.Join(folder, globals.ConfigFileName)
	if err := i.config.Save(filePath); err != nil {
		return err
	}
	return nil
}

func (i *Installer) ConfigFileFolder() (string, error) {
	return config.ConfigFileFolder(i.appID)
}

type InstallStage func() error

func (i *Installer) Stages() []InstallStage {
	return []InstallStage{
		i.StageCreateFolder,
		i.StageCreateUninstallScript,
		i.StageCreateConfig,
	}
}

func (i *Installer) StageCreateConfig() error {
	logging.Debugf("Install: CreateConfig")
	folder, err := i.ConfigFileFolder()
	if err != nil {
		return err
	}
	filePath := filepath.Join(folder, globals.ConfigFileName)
	logging.Debugf("Install: CreateConfig: Save to %s", filePath)
	return i.config.Save(filePath)
}
func (i *Installer) Path(fileName string) string {
	return filepath.Join(i.config.Folder, fileName)
}

const uninstallScriptName = "uninstall"

func (i *Installer) StageCreateUninstallScript() error {
	logging.Debugf("Install: StageCreateUninstallScript")
	var err error
	scriptName := uninstallScriptName + script.Get().Extension()
	i.uninstallScript, err = os.Create(i.Path(scriptName))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(i.uninstallScript, script.Get().Comment("Uninstallation script"))
	return err
}

func (i *Installer) StageCreateFolder() error {
	logging.Debugf("Install: StageCreateFolder %s", i.config.Folder)
	folder := filepath.Join(i.config.Folder, globals.AppName)
	return os.MkdirAll(folder, 0766)
}

func (i *Installer) StageExtractExecutable() error {
	logging.Debugf("Install: StageExtractExecutable")
	return nil
}
