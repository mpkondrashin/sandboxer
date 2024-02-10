// Copyright
package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"sandboxer/pkg/extract"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
)

//go:embed embed/*
var embedFS embed.FS

const wizardLog = globals.Name + "_setup.log"

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func InstallLogFolder() string {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(path)
}

func main() {
	close, err := logging.NewFileLog(InstallLogFolder(), wizardLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewFileLog: %v", err)
		os.Exit(10)
	}
	defer close()
	logging.Infof("Execute Start. Version %s Build %s", globals.Version, globals.Build)

	self, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Executeable: %v", err)
		os.Exit(20)
	}
	logging.Infof("Path: %s", self)
	tempFolder, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getwd: %v", err)
		os.Exit(30)
	}
	logging.Infof("Temp folder: %s", tempFolder)
	if IsWindows() {
		path, err := extract.FileGZ(embedFS, tempFolder, "embed/opengl32.dll.gz")
		logging.LogError(err)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FileGZ: %v", err)
			os.Exit(40)
		}
		logging.Debugf("Extracted: %s", path)
		//defer cleanup()
	}
	installPath, err := extract.FileGZ(embedFS, tempFolder, "embed/install.exe.gz")
	logging.LogError(err)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileGZ: %v", err)
		os.Exit(50)
	}
	logging.Debugf("Execute: %s", installPath)
	cmd := exec.Command(installPath)
	var errb, outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	logging.LogError(err)
	if err != nil {
		logging.Errorf("exit code: %d", cmd.ProcessState.ExitCode())
		if cmd.ProcessState.ExitCode() == 1 {
			logging.Infof("Extracting Open GL")
		}
		logging.Errorf("Error: \"%s\"", errb.String())
		logging.Errorf("Output: \"%s\"", outb.String())
		// Save to file!!!
		return
	}
	logging.Infof("Setup finished")
}
