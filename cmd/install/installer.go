package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	autostart       bool
	uninstallScript *script.Script
}

func NewInstaller(appID string) (*Installer, error) {
	configFolder, err := globals.UserDataFolder()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(configFolder, globals.ConfigFileName)
	logging.Debugf("Configuration path: %s", configPath)

	return &Installer{
		appID:     appID,
		config:    config.New(configPath),
		autostart: true,
	}, nil
}

func (i *Installer) LoadConfig() error {
	err := i.config.Load()
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
	//filePath := filepath.Join(folder, globals.ConfigFileName)
	if err := i.config.Save(); err != nil {
		return err
	}
	return nil
}

func (i *Installer) ConfigFileFolder() (string, error) {
	return globals.UserDataFolder()
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
		{"Uninstall script", i.StageCreateUninstallScript},
		{"Create folders", i.StageCreateFolders},
		{"Generate config", i.StageCreateConfig},
		{"Stop Examen", i.StageStopExamen},
		{"Wait for service to stop", i.StageWaitServiceToStop},
		//{"Uninstall service", i.StageUninstallService},
		{"Extract executables", i.StageExtractExecutable},
		{"Extend Send To menu", i.StageExtendSendTo},
		{"Install service", i.StageAutostart},
		{"Stop runnin " + globals.AppName, i.StageUninstall},
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
	logsFolder, err := globals.LogsFolder()
	if err != nil {
		return err
	}
	folders := []string{
		filepath.Join(i.config.Folder, globals.AppName),
		logsFolder,
	}
	for _, f := range folders {
		logging.Debugf("Install: Create folder \"%s\"", f)
		if err := os.MkdirAll(f, 0766); err != nil {
			return err
		}
		err := i.uninstallScript.AddLine(script.Get().RemoveDir(f))
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
	//filePath := filepath.Join(folder, globals.ConfigFileName)
	//logging.Debugf("Install: CreateConfig: Save to %s", filePath)
	if err := i.config.Save(); err != nil {
		return err
	}
	return i.uninstallScript.AddLine(script.Get().RemoveDir(i.config.GetFilePath()))
}

func (i *Installer) StageStopExamen() error {
	logging.Debugf("Install: StopExamen")
	pidFilePath, err := globals.PidFilePath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(pidFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.Debugf("Install: StopExamen: %s: %v", pidFilePath, err)
			return nil
		}
		return err
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		logging.Debugf("Install: StopExamen: %s: %v", string(data), err)
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		logging.Debugf("Install: StopExamen: FindProcess(%d): %v", pid, err)
		return nil
	}
	if err := proc.Kill(); err != nil {
		logging.Debugf("Install: StopExamen: Kill %d: %v", pid, err)
		return nil
	}
	return nil
}

func (i *Installer) StageWaitServiceToStop() error {
	logging.Debugf("Install: WaitServiceToStop")
	return nil
}

func (i *Installer) StageExtractExecutable() error {
	logging.Debugf("Install: StageExtractExecutable")
	toExtract := []string{
		"embed/examen.exe.gz",
		"embed/submit.exe.gz",
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
	return nil
}

func (i *Installer) StageExtendSendTo() error {
	logging.Debugf("Install: ExtendSendTo")
	appPath := filepath.Join(i.InstallFolder(), "submit.exe")
	linkPath, err := globals.ExtendContextMenu(appPath)
	if err != nil {
		return err
	}
	return i.uninstallScript.AddLine(script.Get().RemoveDir(linkPath))
}

func (i *Installer) StageAutostart() error {
	logging.Debugf("Install: Autostart")
	if !i.autostart {
		logging.Debugf("Install: StageAutostart: Skip")
		return nil
	}
	appPath := filepath.Join(i.InstallFolder(), globals.AppName+".exe")
	linkPath, err := globals.AutoStart(appPath)
	if err != nil {
		return err
	}
	return i.uninstallScript.AddLine(script.Get().RemoveDir(linkPath))
}

func (i *Installer) StageUninstall() error {
	logging.Debugf("Install: Uninstall")
	pidPath, err := globals.PidFilePath()
	if err != nil {
		return err
	}
	err = i.uninstallScript.AddLine(script.Get().RemoveDir(pidPath))
	if err != nil {
		return err
	}
	return i.uninstallScript.AddLine(script.Get().StopProcess(globals.AppName))
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
