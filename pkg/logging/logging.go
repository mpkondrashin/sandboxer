/*
TunnelEffect (c) 2022 by Mikhail Kondrashin (mkondrashin@gmail.com)

logging.go

TunnelEffect logging module.

logging.Debugf()
*/
package logging

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg      sync.WaitGroup
	rw      sync.RWMutex
	loggers []chan LogData
)

func AddLogger(l Logger) {
	rw.Lock()
	defer rw.Unlock()
	c := make(chan LogData, 1000)
	loggers = append(loggers, c)
	wg.Add(1)
	go func() {
		for d := range c {
			l.Write(d)
		}
		wg.Done()
	}()
}

func Close() {
	for _, l := range loggers {
		close(l)
	}
	wg.Wait()
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
	for _, each := range loggers {
		each <- d
	}
}
