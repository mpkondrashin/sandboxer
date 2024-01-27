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

func NewInstaller(appID string) (*Installer, error) {
	configFolder, err := globals.UserDataFolder()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(configFolder, globals.ConfigFileName)
	logging.Debugf("Configuration path: %s", configPath)

	return &Installer{
		appID:  appID,
		config: config.New(configPath),
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
		{"Install service", i.StageAutostart},
		{"Start service", i.StageStart},
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
		i.config.Folder,
		globals.InstallFolder(),
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
	/*
	   os
	   func FindProcess(pid int) (*Process, error)
	   func (p *Process) Kill() error
	   func (p *Process) Signal(sig Signal) error

	*/
	//time.Sleep(1 * time.Second)
	return nil
}

func (i *Installer) StageWaitServiceToStop() error {
	logging.Debugf("Install: WaitServiceToStop")
	return nil
}
func (i *Installer) StageUninstallService() error {
	logging.Debugf("Install: UninstallService")
	s, err := i.config.Service(nil)
	if err != nil {
		return err
	}
	if err := s.Uninstall(); err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("service %s is not installed", globals.SvcName)) {
			return nil
		}
		return err
	}
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

func (i *Installer) StageAutostart() error {
	logging.Debugf("Install: Autostart")
	return nil
}

func (i *Installer) StageStart() error {
	logging.Debugf("Install: Start")
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
