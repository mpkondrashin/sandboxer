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
	// uninstallScript *script.Script
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

/*
	func (i *Installer) __LoadConfig() (ModusOperandi, error) {
		err := i.config.Load()
		if err != nil {
			if os.IsNotExist(err) {
				return moInstall, nil
			}
			return moError, err
		}
		switch semver.Compare(globals.Version, i.config.GetVersion()) {
		case -1:
			return moDowngrade, nil
		case 0:
			return moReinstall, nil
		case 1:
			return moUpgrade, nil
		}
		return moError, errors.New("version error")
	}
*/
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
		//{"Uninstall script", i.StageCreateUninstallScript},
		{"Create folders", i.StageCreateFolders},
		{"Generate config", i.StageCreateConfig},
		{"Stop " + globals.AppName, i.StageStopProgram},
		{"Wait for service to stop", i.StageWaitServiceToStop},
		{"Extract executables", i.StageExtractFiles},
		{"Extend Send To menu", i.StageExtendSendTo},
		{"Add to Start menu", i.StageAddToStartMenu},
		{"Install service", i.StageAutoStart},
		//{"Stop Running " + globals.AppName, i.StageUninstall},
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

//const uninstallScriptName = "uninstall"

/*
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
}*/

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
		/*err := i.uninstallScript.AddLine(script.Get().RemoveDir(f))
		if err != nil {
			return err
		}*/
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
	//i.uninstallScript.AddLine(script.Get().RemoveDir(folder))
	return i.config.Save() /*; err != nil {
		return err
	}
	return i.uninstallScript.AddLine(script.Get().RemoveDir(i.config.GetFilePath()))*/
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
	//if err := i.uninstallScript.AddLine(script.Get().RemoveDir(targetPath)); err != nil {
	//	return err
	//}
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
	//return i.uninstallScript.AddLine(script.Get().RemoveDir(path))
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
	//if err != nil {
	//	return err
	//}
	/*
		appPath := filepath.Join(i.InstallFolder(), xplatform.ExecutableName(globals.Name))
		_, err := xplatform.LinkToStartMenu(false, globals.AppName, globals.AppName, appPath, false)
		if err != nil {
			return err
		}
		scriptPath := i.Path(uninstallScriptName + script.Get().Extension())
		path, err := xplatform.LinkToStartMenu(false, globals.AppName, "Uninstall", scriptPath, true)
		if err != nil {
			return err
		}*/
	//if err := i.uninstallScript.AddLine(script.Get().RemoveDir(path)); err != nil {
	//	return err
	//}
	//return nil
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
	//scriptPath := i.Path(uninstallScriptName + script.Get().Extension())
	//_, err = xplatform.LinkToStartMenu(false, globals.AppName, "Uninstall", scriptPath, true)
	//if err != nil {
	//		return "", err
	//	}
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
	//return i.uninstallScript.AddLine(script.Get().RemoveDir(path))
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

/*func (i *Installer) StageUninstall() error {
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
}*/
/*scriptName := uninstallScriptName + script.Get().Extension()
var scriptPath string
if xplatform.IsWindows() {
	scriptPath = i.Path(scriptName)
} else {
	scriptPath, err = xplatform.ExecutablePath(xplatform.InstallFolder(), globals.AppName, scriptName)
	if err != nil {
		return err
	}
}
/*if err := os.Rename(i.uninstallScript.FilePath, scriptPath); err != nil {
	return err
}*/
// remove folder
//return nil
//}

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
