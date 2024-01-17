package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"examen/pkg/config"
	"examen/pkg/extract"
	"examen/pkg/globals"
	"examen/pkg/logging"
	"examen/pkg/script"
)

//go:embed embed/*
var embedFS embed.FS

type Installer struct {
	appID           string
	config          *config.Configuration
	uninstallScript *script.Script
}

func NewInstaller(appID string) *Installer {
	return &Installer{
		appID:  appID,
		config: config.New(),
	}
}

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

func (i *Installer) SaveConfig() error {
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

func (i *Installer) Path(fileName string) string {
	return filepath.Join(i.config.Folder, globals.AppFolderName, fileName)
}

func (i *Installer) InstallFolder() string {
	return filepath.Join(i.config.Folder, globals.AppFolderName)
}

type InstallStage func() error

func (i *Installer) Stages() []InstallStage {
	return []InstallStage{
		i.StageCreateUninstallScript,
		i.StageCreateFolder,
		i.StageCreateConfig,
		i.StageExtractExecutable,
	}
}

const uninstallScriptName = "uninstall"

func (i *Installer) StageCreateUninstallScript() error {
	logging.Debugf("Install: StageCreateUninstallScript")
	scriptName := uninstallScriptName + script.Get().Extension()
	logging.Debugf("Install:  uninstall script name: %s", scriptName)
	i.uninstallScript = script.New(i.Path(scriptName), script.Get().Comment("Uninstallation script"))
	return nil
}

func (i *Installer) StageCreateFolder() error {
	logging.Debugf("Install: StageCreateFolder")
	folder := filepath.Join(i.config.Folder, globals.AppFolderName)
	logging.Debugf("Install: Create folder \"%s\"", folder)
	if err := os.MkdirAll(folder, 0766); err != nil {
		return err
	}
	return i.uninstallScript.AddLine(script.Get().RemoveDir(folder))
}

func (i *Installer) StageCreateConfig() error {
	logging.Debugf("Install: CreateConfig")
	folder, err := i.ConfigFileFolder()
	if err != nil {
		return err
	}
	logging.Debugf("Create folder \"%s\"", folder)
	if err := os.MkdirAll(folder, 0500); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	i.uninstallScript.AddLine(script.Get().RemoveDir(folder))
	filePath := filepath.Join(folder, globals.ConfigFileName)
	logging.Debugf("Install: CreateConfig: Save to %s", filePath)
	if err := i.config.Save(filePath); err != nil {
		return err
	}
	return i.uninstallScript.AddLine(script.Get().RemoveDir(filePath))
}

func (i *Installer) StageExtractExecutable() error {
	logging.Debugf("Install: StageExtractExecutable")
	if IsWindows() {
		path, err := extract.FileGZ(embedFS, i.InstallFolder(), "embed/opengl32.dll.gz")
		if err != nil {
			return err
		}
		logging.Debugf("Extracted: %s", path)
	}
	examensvcPath, err := extract.FileGZ(embedFS, i.InstallFolder(), "embed/examensvc.exe.gz")
	if err != nil {
		return err
	}
	logging.Debugf("Extracted: %s", examensvcPath)
	examenPath, err := extract.FileGZ(embedFS, i.InstallFolder(), "embed/examen.exe.gz")
	if err != nil {
		return err
	}
	logging.Debugf("Extracted: %s", examenPath)
	return nil
}

func (i *Installer) StageExtractPericulosum() error {
	logging.Debugf("Install: StageExtractPericulosum")
	///
	return nil
}

func (i *Installer) StageRightClickExtension() error {
	logging.Debugf("Install: StageRightClickExtension")
	// Computer/HCU/Directory/Background//shell
	// + New
	// Examen
	// New key

	// command=app path %1
	////
	return nil
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
