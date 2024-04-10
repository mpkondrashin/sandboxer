/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

installer.go

Installer struct
*/
package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"sandboxer/pkg/config"
	"sandboxer/pkg/extract"
	"sandboxer/pkg/fifo"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/xplatform"
)

//go:embed embed/*.tar.gz embed/LICENSE
var embedFS embed.FS

type Installer struct {
	appID     string
	config    *config.Configuration
	autostart bool
}

func NewInstaller(appID string) (*Installer, error) {
	configFolder, err := xplatform.UserDataFolder(appID)
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

func (i *Installer) SaveConfig() error {
	logging.Debugf("Save config")
	folder, err := i.ConfigFileFolder()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(folder, 0700); err != nil {
		return err
	}
	if err := i.config.Save(); err != nil {
		return err
	}
	return nil
}

func (i *Installer) ConfigFileFolder() (string, error) {
	return xplatform.UserDataFolder(globals.AppID)
}

func (i *Installer) Path(fileName string) string {
	return filepath.Join(i.config.GetFolder(), globals.AppFolderName, fileName)
}

func (i *Installer) InstallFolder() string {
	if xplatform.IsWindows() {
		return filepath.Join(i.config.GetFolder(), globals.AppFolderName)
	} else { // It is macOS
		return filepath.Join(i.config.GetFolder(), globals.AppName+".app")
	}
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
		{"Create folders", i.StageCreateFolders},
		{"Generate config", i.StageCreateConfig},
		{"Stop " + globals.AppName, i.StageStopProgram},
		{"Wait for service to stop", i.StageWaitServiceToStop},
		{"Extract executables", i.StageExtractFiles},
		{"Extend Send To menu", i.StageExtendSendTo},
		{"Add to Start menu", i.StageAddToStartMenu},
		{"Install service", i.StageAutoStart},
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

type UninstallStage interface {
	Name() string
	Execute() error
}

func (i *Installer) UninstallStages() (stages []UninstallStage) {
	stages = []UninstallStage{
		NewUninstallStageStopProcess(),
		NewUninstallStageUnregister(i.config.DDAn),
		NewUninstallStageDelete("Program", i.InstallFolder()),
	}
	configFolder, err := i.ConfigFileFolder()
	if err == nil {
		stages = append(stages, NewUninstallStageDelete("User Data", configFolder))
	}
	autoStartPath, err := i.AutoStart(false)
	if err == nil {
		stages = append(stages, NewUninstallStageDelete("Autostart", autoStartPath))
	}
	sendTo, err := i.ExtendSendTo(true)
	if err == nil {
		stages = append(stages, NewUninstallStageDelete("Right Click", sendTo))
	}
	if xplatform.IsWindows() {
		path, err := i.AddToStartMenu(true)
		if err == nil {
			stages = append(stages, NewUninstallStageDelete("Start Menu", path))
		}
	}
	stages = append(stages, NewUninstallStageDone())
	return
}

func (i *Installer) StageCreateFolders() error {
	logging.Debugf("Install: StageCreateFolders")
	logsFolder, err := globals.LogsFolder()
	if err != nil {
		return err
	}
	tasksFolder, err := globals.TasksFolder()
	if err != nil {
		return err
	}
	folders := []string{
		logsFolder,
		tasksFolder,
	}
	if xplatform.IsWindows() {
		folders = append(folders, i.InstallFolder())
	}
	for _, f := range folders {
		logging.Debugf("Install: Create folder \"%s\"", f)
		if err := os.MkdirAll(f, 0755); err != nil {
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
	return i.config.Save()
}

func (i *Installer) StageStopProgram() error {
	logging.Debugf("Install: Stop " + globals.AppName)
	return NewUninstallStageStopProcess().Execute()
}

func (i *Installer) StageWaitServiceToStop() error {
	logging.Debugf("Install: WaitServiceToStop")
	for i := 0; i < 10; i++ {
		fifoWriter, err := fifo.NewWriter()
		if err != nil && fifo.IsDown(err) {
			return nil
		}
		fifoWriter.Close()
		logging.Debugf("Wait for FIFO")
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("stop Submissions and run setup again")
}

func (i *Installer) StageExtractFiles() error {
	logging.Debugf("Install: StageExtractExecutable")
	path := "embed/" + globals.Name + ".tar.gz"
	targetPath := i.InstallFolder()
	if !xplatform.IsWindows() {
		targetPath = filepath.Dir(targetPath)
	}
	if err := extract.Untar(embedFS, targetPath, path); err != nil {
		return err
	}
	logging.Debugf("Extracted: %s", targetPath)
	return nil
}

func (i *Installer) StageExtractPericulosum() error {
	logging.Debugf("Install: StageExtractPericulosum")
	return nil
}

func (i *Installer) StageExtendSendTo() error {
	logging.Debugf("Install: ExtendSendTo")
	_, err := i.ExtendSendTo(false)
	if err != nil {
		return err
	}
	return err
}

func (i *Installer) ExtendSendTo(dryRun bool) (string, error) {
	if xplatform.IsWindows() {
		appPath := filepath.Join(i.InstallFolder(), "submit.exe")
		return xplatform.ExtendContextMenu(dryRun, globals.AppName, appPath)
	} else { // macOS
		src := "embed/" + globals.Name + "_submit.tar.gz"
		folder := filepath.Join(os.Getenv("HOME"), "/Library/Services")
		workflowPath := filepath.Join(folder, globals.AppName+".workflow")
		if dryRun {
			return workflowPath, nil
		}
		err := extract.Untar(embedFS, folder, src)
		if err != nil {
			return "", err
		}
		return workflowPath, nil
	}
}

func (i *Installer) StageAddToStartMenu() error {
	logging.Infof("Install: AddToStartMenu")
	if !xplatform.IsWindows() {
		logging.Infof("AddToStartMenu. This is not Windows. Skip")
		return nil
	}
	_, err := i.AddToStartMenu(false)
	return err
}

func (i *Installer) AddToStartMenu(dryRun bool) (string, error) {
	appPath := filepath.Join(i.InstallFolder(), xplatform.ExecutableName(globals.Name))
	path, err := xplatform.LinkToStartMenu(dryRun, globals.AppName, globals.AppName, appPath, false)
	if err != nil {
		return "", err
	}
	if dryRun {
		return filepath.Dir(path), nil
	}
	return filepath.Dir(path), nil
}

func (i *Installer) StageAutoStart() error {
	logging.Debugf("Install: Autostart")
	if !i.autostart {
		logging.Debugf("Install: StageAutostart: Skip")
		return nil
	}
	_, err := i.AutoStart(false)
	if err != nil {
		return fmt.Errorf("AutoStart: %w", err)
	}
	return nil
}

func (i *Installer) AutoStart(dryRun bool) (string, error) {
	appPath, err := xplatform.ExecutablePath(i.config.GetFolder(), globals.AppName, globals.Name)
	if err != nil {
		return "", fmt.Errorf("ExecutablePath: %w", err)
	}
	path, err := xplatform.AutoStart(false, globals.AppID, appPath)
	if err != nil {
		return "", fmt.Errorf("AutoStart: %w", err)
	}
	return path, err
}
