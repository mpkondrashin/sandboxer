/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

folders.go

Create new file log
*/
package logging

import (
	"os"
	"path/filepath"
)

func NewFileLog(folder, fileName string) (func(), error) {
	SetLevel(DEBUG)
	logFilePath := filepath.Join(folder, fileName)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	SetLogger(NewFileLogger(file))
	return func() {
		file.Close()
	}, nil
}
