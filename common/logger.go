package common

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
	level  LogLevel
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelStrings = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// NewLogger creates a new logger that writes to the given writer (e.g. os.Stdout or a file)
// and filters messages below the provided log level.
func NewLogger(out io.Writer, level LogLevel) *Logger {
	return &Logger{
		logger: log.New(out, "", 0), // no prefix, we add timestamp etc ourselves
		level:  level,
	}
}

// formatPrefix formats the timestamp and level prefix.
func (l *Logger) formatPrefix(level LogLevel) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	levelStr := levelStrings[level]
	return "[" + t + "] [" + levelStr + "] "
}

func (l *Logger) log(level LogLevel, v ...interface{}) {
	if level < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	prefix := l.formatPrefix(level)
	l.logger.Println(prefix + fmt.Sprint(v...))
}

func (l *Logger) logf(level LogLevel, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	prefix := l.formatPrefix(level)
	l.logger.Println(prefix + fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.log(DEBUG, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(DEBUG, format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.log(INFO, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(INFO, format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(WARN, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logf(WARN, format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.log(ERROR, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(ERROR, format, v...)
}
