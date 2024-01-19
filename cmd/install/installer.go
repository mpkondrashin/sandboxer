package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"examen/pkg/config"
	"examen/pkg/extract"
	"examen/pkg/globals"
	"examen/pkg/logging"
	"examen/pkg/script"
)

//go:embed embed/*.gz
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

func (i *Installer) Run(name string, args ...string) error {
	logging.Debugf("Run %s %s", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...) // #nosec
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	stdout := outBuf.String()
	stderr := errBuf.String()
	if err != nil {
		return fmt.Errorf("error %w: %s, %s", err, stderr, stdout)
	}
	return nil
}

type InstallStage struct {
	Name string
	Run  func() error
}

func (i *Installer) Stages() []InstallStage {
	return []InstallStage{
		{"Uninstall Script", i.StageCreateUninstallScript},
		{"Create Folders", i.StageCreateFolders},
		{"Generate Config", i.StageCreateConfig},
		{"Stop Service", i.StageStopService},
		{"Uninstall Service", i.StageUninstallService},
		{"Extract Executables", i.StageExtractExecutable},
		{"Install Service", i.StageInstallService},
	}
}

func (i *Installer) Install(callback func(name string) error) error {
	for _, stage := range i.Stages() {
		logging.Debugf("Install Stage %s", stage.Name)
		if err := callback(stage.Name); err != nil {
			return err
		}
		if err := stage.Run(); err != nil {
			return err
		}
	}
	return nil
}

const uninstallScriptName = "uninstall"

func (i *Installer) StageCreateUninstallScript() error {
	logging.Debugf("Install: StageCreateUninstallScript")
	scriptName := uninstallScriptName + script.Get().Extension()
	logging.Debugf("Install:  uninstall script name: %s", scriptName)
	i.uninstallScript = script.New(i.Path(scriptName), script.Get().Comment("Uninstallation script"))
	return nil
}

func (i *Installer) StageCreateFolders() error {
	logging.Debugf("Install: StageCreateFolders")
	folders := []string{
		"",
		"logs",
	}
	for _, f := range folders {
		folder := filepath.Join(i.config.Folder, globals.AppFolderName, f)
		logging.Debugf("Install: Create folder \"%s\"", folder)
		if err := os.MkdirAll(folder, 0766); err != nil {
			return err
		}
		err := i.uninstallScript.AddLine(script.Get().RemoveDir(folder))
		if err != nil {
			return err
		}
	}
	return nil
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

func (i *Installer) StageStopService() error {
	logging.Debugf("Install: StopService")
	s, err := i.config.Service(nil)
	if err != nil {
		return err
	}
	if err := s.Stop(); err != nil {
		logging.Debugf("err: %v, err = %T", err, err)
		if !strings.Contains(err.Error(), "The service has not been started") {
			return err
		}
	}
	return nil
}

func (i *Installer) StageUninstallService() error {
	logging.Debugf("Install: UninstallService")
	s, err := i.config.Service(nil)
	if err != nil {
		return err
	}
	return s.Uninstall()
}

func (i *Installer) StageExtractExecutable() error {
	logging.Debugf("Install: StageExtractExecutable")
	toExtract := []string{
		"embed/examensvc.exe.gz",
		"embed/examen.exe.gz",
	}
	if IsWindows() {
		toExtract = append(toExtract, "embed/opengl32.dll.gz")
	}
	for _, path := range toExtract {
		newPath, err := extract.FileGZ(embedFS, i.InstallFolder(), path)
		if err != nil {
			return err
		}
		logging.Debugf("Extracted: %s", newPath)
		if err := i.uninstallScript.AddLine(script.Get().RemoveDir(newPath)); err != nil {
			return err
		}
	}
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

func (i *Installer) StageInstallService() error {
	s, err := i.config.Service(nil)
	if err != nil {
		return err
	}
	if err := s.Install(); err != nil {
		return err
	}
	return i.uninstallScript.AddLine(script.Get().UninstallService(globals.SvcName))
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
