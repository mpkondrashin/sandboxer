package main

import (
	"bytes"
	"examen/pkg/logging"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func setupLogging(logFileName string) func() {
	logFolder, err := os.Executable()
	if err != nil {
		//logging.Errorf("os.Executable: %v", err)
		panic(err)
	}
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
	//redirectStderr(file)
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

const (
	examenSetupWizardLog   = "examen_setup_wizard.log"
	examenExecuteWizardLog = "examen_execute_wizard.log"
)

func runWizard() {
	close := setupLogging(examenSetupWizardLog)
	defer close()
	/*defer func() {
		if err := recover(); err != nil {
			logging.Criticalf("panic: %v", err)
		}
	}()*/
	logging.Infof("Setup Start")
	//if IsWindows() {
	//	cleanup := extractOpenGL()
	//	defer cleanup()
	//}
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
	close := setupLogging(examenExecuteWizardLog)
	defer close()
	logging.Infof("Execute Start")
	self, err := os.Executable()
	if err != nil {
		panic(err)
	}
	logging.Infof("Path: %s", self)
	for i := 0; i < 2; i++ {
		cmd := exec.Command(self, runGUIparameter)
		var errb bytes.Buffer
		//cmd.Stdout = &outb
		cmd.Stderr = &errb
		err = cmd.Run()
		logging.LogError(err)
		if err != nil {
			if cmd.ProcessState.ExitCode() == 1 {
				logging.Infof("Extracting Open GL")
				extractOpenGL()
				continue
			}
			logging.Errorf("Error: %v", errb)
			return
		}
		break
		//fmt.Println("err", err, "/err")
		//fmt.Println("EXITCODE", cmd.ProcessState.ExitCode(), "EXITCODE")
		//fmt.Println("Stdout", outb.String(), "/Stdout")
		//fmt.Println("Stderr", errb.String(), "/Stderr")
	}
	logging.Infof("Execute Stop")
}
func main() {
	if len(os.Args) == 2 && os.Args[1] == runGUIparameter {
		runWizard()
		return
	}
	executeWizard()
}
