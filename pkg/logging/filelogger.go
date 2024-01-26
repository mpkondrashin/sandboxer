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
