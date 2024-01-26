package logging

import (
	"os"
	"path/filepath"
)

func NewFileLog(folder, fileName string) func() {
	SetLevel(DEBUG)
	logFilePath := filepath.Join(folder, fileName)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	SetLogger(NewFileLogger(file))
	return func() {
		file.Close()
	}
}
