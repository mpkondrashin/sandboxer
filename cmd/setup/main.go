package main

import (
	"examen/pkg/logging"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func setupLogging(logFileName string) func() {
	path, err := os.Executable()
	if err != nil {
		//logging.Errorf("os.Executable: %v", err)
		panic(err)
	}
	logFolder := filepath.Dir(path)
	//logFolder := "." //os.TempDir()
	/*	errFileName := "examen_stderr.log"
		errFilePath := filepath.Join(logFolder, errFileName)
		errFile, err := os.OpenFile(errFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(errFile, "%v Started\n", time.Now())
	*/
	//redirectStderr(errFile)

	//installLogFileName := "setup.log"
	logFilePath := filepath.Join(logFolder, logFileName)
	//fmt.Println(logFilePath)
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
	configFileName = "examen.yaml" // remove - use fyne
	appID          = "com.github.mpkondrashin.examen"
	setupWizardLog = "examen_setup_wizard.log"
	openGLdll_gz   = "opengl32.dll.gz"
)

func extractOpenGL() func() {
	path, err := os.Executable()
	if err != nil {
		logging.Errorf("os.Executable: %v", err)
		panic(err)
	}
	folder := filepath.Dir(path)
	filePath, err := extractEmbeddedGZ(folder, openGLdll_gz)
	if err != nil {
		panic(fmt.Errorf("extract %s: %w", openGLdll_gz, err))
	}
	return func() {
		os.Remove(filePath)
	}
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func main() {
	close := setupLogging(setupWizardLog)
	defer func() {
		logging.Debugf("Close log file")
		close()
	}()
	defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()
	logging.Infof("Start")
	logging.Debugf("OS: %s (%s)", runtime.GOOS, runtime.GOARCH)
	if IsWindows() {
		cleanup := extractOpenGL()
		defer func() {
			logging.Debugf("Remove OpenGL DLL")
			cleanup()
		}()
	}
	capturesFolder := ""
	if len(os.Args) == 3 && os.Args[1] == "--capture" {
		capturesFolder = os.Args[2]
	}
	logging.Infof("Starting Wizard")
	c := NewNSHIControl(capturesFolder)
	c.Run()
	logging.Infof("Setup finished")
}
