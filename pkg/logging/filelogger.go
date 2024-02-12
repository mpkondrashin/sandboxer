/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

logger.go

File logger
*/
package logging

import (
	"io"
)

type FileLogger struct {
	w io.Writer
}

var _ Logger = &FileLogger{}

func NewFileLogger(w io.Writer) *FileLogger {
	return &FileLogger{
		w: w,
	}
}

func (fl *FileLogger) Write(data LogData) {
	_, _ = fl.w.Write([]byte(data.String()))
	_, _ = fl.w.Write([]byte("\n"))
}
