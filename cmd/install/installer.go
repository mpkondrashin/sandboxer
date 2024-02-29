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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"sandboxer/pkg/config"
	"sandboxer/pkg/extract"
	"sandboxer/pkg/fifo"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/script"
	"sandboxer/pkg/xplatform"
)

//go:embed embed/*.tar.gz
var embedFS embed.FS

type Installer struct {
	appID           string
	config          *config.Configuration
	autostart       bool
	uninstallScript *script.Script
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
	return xplatform.UserDataFolder(globals.AppID)
}

func (i *Installer) Path(fileName string) string {
	return filepath.Join(i.config.Folder, globals.AppFolderName, fileName)
}

func (i *Installer) InstallFolder() string {
	if xplatform.IsWindows() {
		return filepath.Join(i.config.Folder, globals.AppFolderName)
	} else {
		return i.config.Folder
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
		{"Uninstall script", i.StageCreateUninstallScript},
		{"Create folders", i.StageCreateFolders},
		{"Generate config", i.StageCreateConfig},
		{"Stop " + globals.AppName, i.StageStopProgram},
		{"Wait for service to stop", i.StageWaitServiceToStop},
		//{"Uninstall service", i.StageUninstallService},
		{"Extract executables", i.StageExtractFiles},
		{"Extend Send To menu", i.StageExtendSendTo},
		{"Install service", i.StageAutoStart},
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
	logging.Debugf("Install: uninstall script name: %s", scriptName)
	temp, err := os.MkdirTemp(os.TempDir(), globals.Name+"-uninstall-*")
	if err != nil {
		return fmt.Errorf("MkdirTemp: %w", err)
	}
	uninstallScriptPath := filepath.Join(temp, scriptName)
	logging.Debugf("Install: uninstall script path: %s", uninstallScriptPath)
	i.uninstallScript = script.New(uninstallScriptPath, script.Get().Comment("Uninstallation script"))
	return nil
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

func (i *Installer) StageStopProgram() error {
	logging.Debugf("Install: Stop " + globals.AppName)
	pidFilePath, err := globals.PidFilePath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(pidFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.Debugf("Install: Stop "+globals.AppName+": %s: %v", pidFilePath, err)
			return nil
		}
		return err
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		logging.Debugf("Install: Stop"+globals.AppName+": %s: %v", string(data), err)
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		logging.Debugf("Install: Stop"+globals.AppName+": FindProcess(%d): %v", pid, err)
		return nil
	}
	if err := proc.Kill(); err != nil {
		logging.Debugf("Install: Stop"+globals.AppName+": Kill %d: %v", pid, err)
		return nil
	}
	return nil
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
	if err := extract.Untar(embedFS, i.InstallFolder(), path); err != nil {
		return err
	}
	targetPath := i.InstallFolder()
	if !xplatform.IsWindows() {
		targetPath = filepath.Join(i.InstallFolder(), globals.AppName+".app")
	}
	if err := i.uninstallScript.AddLine(script.Get().RemoveDir(targetPath)); err != nil {
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
	var path string
	var err error
	if runtime.GOOS == "windows" {
		appPath := filepath.Join(i.InstallFolder(), "submit.exe")
		path, err = xplatform.ExtendContextMenu(globals.AppName, appPath)
		if err != nil {
			return err
		}
	} else {
		src := "embed/" + globals.Name + "_submit.tar.gz"
		folder := filepath.Join(os.Getenv("HOME"), "/Library/Services")
		err := extract.Untar(embedFS, folder, src)
		if err != nil {
			return err
		}
		path = filepath.Join(folder, globals.AppName+".workflow")
	}
	return i.uninstallScript.AddLine(script.Get().RemoveDir(path))
}

func (i *Installer) StageAutoStart() error {
	logging.Debugf("Install: Autostart")
	if !i.autostart {
		logging.Debugf("Install: StageAutostart: Skip")
		return nil
	}
	var appPath string
	appPath, err := xplatform.ExecutablePath(i.config.Folder, globals.AppName, globals.Name)
	if err != nil {
		return fmt.Errorf("ExecutablePath: %w", err)
	}
	path, err := xplatform.AutoStart(globals.AppID, appPath)
	if err != nil {
		return fmt.Errorf("AutoStart: %w", err)
	}
	return i.uninstallScript.AddLine(script.Get().RemoveDir(path))
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
	if err := i.uninstallScript.AddLine(script.Get().StopProcess(pidPath)); err != nil {
		return err
	}
	scriptName := uninstallScriptName + script.Get().Extension()
	var scriptPath string
	if xplatform.IsWindows() {
		scriptPath = i.Path(scriptName)
	} else {
		scriptPath, err = xplatform.ExecutablePath(xplatform.InstallFolder(), globals.AppName, scriptName)
		if err != nil {
			return err
		}
	}
	if err := os.Rename(i.uninstallScript.FilePath, scriptPath); err != nil {
		return err
	}
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
