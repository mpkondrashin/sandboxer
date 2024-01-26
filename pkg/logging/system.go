package logging

import (
	"errors"
	"fmt"
	"io/fs"
	"runtime"
	"strings"
)

var loggingLevel = INFO

// var name = ""
var ErrUnknownLogLevel = errors.New("unknow log level")

// SetLevelStr - set logging level using string level representation.
func SetLevelStr(level string) error {
	for i := DEBUG; i <= CRITICAL; i++ {
		if strings.EqualFold(LevelName[i], level) {
			SetLevel(i)

			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknownLogLevel, level)
}

// SetLevel - set logging level.
func SetLevel(level int) {
	rw.Lock()
	defer rw.Unlock()
	if level < DEBUG || level > CRITICAL {
		return
	}
	loggingLevel = level
}

/*
	func SetName(n string) {
		rw.Lock()
		defer rw.Unlock()
		name = " " + n
	}
*/
type System struct {
	Time     string
	Severity string
	Thread   string
	Message  string
	Path     string
}

var _ LogData = &System{}

func (s *System) String() string {
	return fmt.Sprintf("%s %s %s %s %s",
		s.Time /*name,*/, s.Severity, s.Thread, s.Message, s.Path)
}

// Logging levels.
const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

// LevelName - loging levels as strings.
var LevelName = []string{
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"CRITICAL",
}

// Debugf - log debug level message.
func Debugf(format string, v ...interface{}) {
	logItf(DEBUG, format, v...)
}

// Infof - log info level message.
func Infof(format string, v ...interface{}) {
	logItf(INFO, format, v...)
}

// Warningf - log warning level message.
func Warningf(format string, v ...interface{}) {
	logItf(WARNING, format, v...)
}

// Errorf - log error level message.
func Errorf(format string, v ...interface{}) {
	logItf(ERROR, format, v...)
}

// Criticalf - log critical level message.
func Criticalf(format string, v ...interface{}) {
	logItf(CRITICAL, format, v...)
}

// LogError - log error if err is not nil.
func LogError(err error) {
	if err == nil {
		return
	}
	msg := ""
	if errors.Is(err, fs.ErrInvalid) {
		msg = "ErrInvalid"
	}
	if errors.Is(err, fs.ErrPermission) {
		msg = "ErrPermission"
	}
	if errors.Is(err, fs.ErrExist) {
		msg = "ErrExist"
	}
	if errors.Is(err, fs.ErrNotExist) {
		msg = "ErrNotExist"
	}
	if errors.Is(err, fs.ErrClosed) {
		msg = "ErrClosed"
	}
	if msg != "" {
		msg += ": " + err.Error()
	} else {
		msg = err.Error()
	}
	logItf(ERROR, msg)
}

func logItf(severity int, format string, v ...interface{}) {
	rw.Lock()
	defer rw.Unlock()
	if severity < loggingLevel {
		return
	}
	message := fmt.Sprintf(format, v...)
	message = strings.ReplaceAll(message, "\n", "\\n")
	message = strings.ReplaceAll(message, "\r", "\\r")

	system := System{
		Time:     timeString(),
		Severity: LevelName[severity],
		Thread:   fmt.Sprintf("G%04s", GoRoutineNumber()),
		Message:  message,
		Path:     callPath(skipLevels),
	}
	logData(&system)
}

const skipLevels = 3

func callPath(skip int) string {
	var sb strings.Builder
	separator := ""
	for c := skip; c < 100; c++ {
		pc, _, line, ok := runtime.Caller(c)
		if !ok {
			break
		}

		name := runtime.FuncForPC(pc).Name()
		if name == "main.main" || strings.HasPrefix(name, "runtime") {
			continue
		}

		p := strings.LastIndex(name, "/")
		if p != -1 {
			// Remove path to module
			name = name[p+1:]
		}

		sb.WriteString(fmt.Sprintf("%s%s[%d]", separator, name, line))
		separator = "<"
	}
	return sb.String()
}

func GoRoutineNumber() string {
	buf := make([]byte, 20)
	runtime.Stack(buf, false)
	s := string(buf)
	space := strings.Index(s[10:], " ")
	if space == -1 {
		return "?"
	}
	return s[10 : 10+space]
}
