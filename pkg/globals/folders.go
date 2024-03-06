/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

folders.go

Various folders that used through whole project
*/
package globals

import (
	"os"
	"path/filepath"

	"sandboxer/pkg/logging"
	"sandboxer/pkg/xplatform"

	"github.com/virtuald/go-paniclog"
)

const (
	tasksFolder = "tasks"
	logsFolder  = "logs"
)

func ConfigurationFilePath() (string, error) {
	folder, err := xplatform.UserDataFolder(AppID)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, ConfigFileName), nil
}

func LogsFolder() (string, error) {
	folder, err := xplatform.UserDataFolder(AppID)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, logsFolder), nil
}

func TasksFolder() (string, error) {
	folder, err := xplatform.UserDataFolder(AppID)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, tasksFolder), nil
}

func AnalyzerClientUUIDFilePath() (string, error) {
	folder, err := xplatform.UserDataFolder(AppID)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, AnalyzerClientUUID), nil
}

func PidFilePath() (string, error) {
	folder, err := xplatform.UserDataFolder(AppID)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, Name+".pid"), nil
}

func SetupLogging(logFileName string) (func(), error) {
	logging.SetLevel(logging.DEBUG)
	logFolder, err := LogsFolder()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(logFolder, 0755); err != nil {
		return nil, err
	}
	logging.SetLevel(logging.DEBUG)
	file, err := logging.OpenRotated(logFolder, logFileName, 0644, MaxLogFileSize, LogsKeep)
	if err != nil {
		return nil, err
	}
	paniclog.RedirectStderr(file.File)
	logging.SetLogger(logging.NewFileLogger(file))
	return func() {
		logging.Infof("Close Logging")
		file.Close()
	}, nil
}
