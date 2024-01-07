package main

import (
	"examen/pkg/logging"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func setupLogging() func() {
	logFolder := "tmp" //os.TempDir()

	/*	errFileName := "examen_stderr.log"
		errFilePath := filepath.Join(logFolder, errFileName)
		errFile, err := os.OpenFile(errFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(errFile, "%v Started\n", time.Now())
	*/
	//redirectStderr(errFile)

	installLogFileName := "setup.log"
	logFilePath := filepath.Join(logFolder, installLogFileName)
	fmt.Println(logFilePath)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println(logFilePath)
	redirectStderr(file)
	logger := logging.NewFileLogger(file)
	logging.AddLogger(logger)
	logging.SetLevel(logging.DEBUG)
	return func() {
		logging.Close()
		file.Close()
		//	errFile.Close()
	}
}

const (
	appName        = "Examen"
	configFileName = "examen.yaml"
	appID          = "com.github.mpkondrashin.examen"
)

func extractOpenGL() func() {
	path, err := os.Executable()
	if err != nil {
		logging.Errorf("os.Executable: %v", err)
		panic(err)
	}
	folder := filepath.Dir(path)
	filePath := extractEmbeddedGZ(folder, "opengl32.dll.gz")
	return func() {
		os.Remove(filePath)
	}
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func runWizard() {
	close := setupLogging()
	defer close()
	/*defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()*/
	logging.Infof("Setup Start")
	if IsWindows() {
		cleanup := extractOpenGL()
		defer cleanup()
	}
	capturesFolder := ""
	if len(os.Args) == 3 && os.Args[1] == "--capture" {
		capturesFolder = os.Args[2]
	}
	c := NewNSHIControl(capturesFolder)
	c.Run()
	logging.Infof("Setup finished")
}

const runGUIparameter = "gui"

func executeWizard() {
	self, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cmdOutput := exec.Command(self, runGUIparameter)
	if err := cmdOutput.Run(); err != nil {
		panic(err)
	}
	fmt.Println(cmdOutput.ProcessState.ExitCode())
}
func main() {
	if len(os.Args) == 2 && os.Args[1] == runGUIparameter {
		runWizard()
		return
	}
	// Inserted code
	executeWizard()

}
