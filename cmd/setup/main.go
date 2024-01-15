// Copyright
package main

import (
	"bytes"
	"embed"
	"os"
	"os/exec"
	"runtime"

	"examen/pkg/extract"
	"examen/pkg/logging"
)

//go:embed embed/*
var embedFS embed.FS

const examenExecuteWizardLog = "examen_setup.log"

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func main() {
	close := logging.NewFileLog(logging.InstallLogFolder(), examenExecuteWizardLog)
	defer close()
	logging.Infof("Execute Start")
	self, err := os.Executable()
	if err != nil {
		panic(err)
	}
	logging.Infof("Path: %s", self)
	tempFolder := "."
	if IsWindows() {
		path, err := extract.FileGZ(embedFS, tempFolder, "embed/opengl32.dll.gz")
		logging.LogError(err)
		if err != nil {
			panic(err)
		}
		logging.Debugf("Extracted: %s", path)
		//defer cleanup()
	}
	installPath, err := extract.FileGZ(embedFS, tempFolder, "embed/install.exe.gz")
	logging.LogError(err)
	if err != nil {
		panic(err)
	}
	logging.Debugf("Extracted: %s", installPath)

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
