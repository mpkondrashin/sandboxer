package logging

import (
	"os"
	"path/filepath"
)

func InstallLogFolder() string {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(path)
}

func NewFileLog(folder, fileName string) func() {
	SetLevel(DEBUG)
	logFilePath := filepath.Join(folder, fileName)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logger := NewFileLogger(file)
	AddLogger(logger)
	return func() {
		Close()
		file.Close()
	}
}
