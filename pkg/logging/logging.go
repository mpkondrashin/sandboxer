/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

logging.go

Logging functions
*/
package logging

import (
	"fmt"
	"sync"
	"time"
)

var (
	//wg     sync.WaitGroup
	rw     sync.RWMutex
	logger Logger
)

func SetLogger(l Logger) {
	rw.Lock()
	defer rw.Unlock()
	logger = l
}

func SetTimeFormat(format string) {
	rw.Lock()
	defer rw.Unlock()
	timeFormat = format
}

type LogData interface {
	fmt.Stringer
}

type Logger interface {
	Write(LogData)
}

var timeFormat = "2006-01-02 15:04:05.000000"

func timeString() string {
	return time.Now().Format(timeFormat)
}

func logData(d LogData) {
	logger.Write(d)
}
