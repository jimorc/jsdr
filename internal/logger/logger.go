package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type LoggingLevel uint8

// Logging levels
const (
	_                  = iota
	Fatal LoggingLevel = iota
	Error
	Info
	Debug
)

var LevelsAsStrings [5]string = [5]string{"Undefined", "Fatal", "Error", "Info", "Debug"}

type logMessage struct {
	level   LoggingLevel
	message string
}

// Logger is a simple logger. It provides a few functions and methods to log information
// to any io.StringWriter.
type Logger struct {
	writer    io.StringWriter
	file      *os.File
	waitGroup sync.WaitGroup
	level     LoggingLevel
	logCh     chan logMessage
}

// New creates a new Logger with a max logging level of Info.
func New(writer io.StringWriter) *Logger {
	l := &Logger{writer: writer}
	l.waitGroup.Add(1)
	l.logCh = make(chan logMessage, 100)
	go l.outputMessages()
	l.SetMaxLevel(Info)
	return l
}

// NewFileLogger creates a new Logger that logs to a file.
func NewFileLogger(f string) (*Logger, error) {
	var file *os.File
	var err error
	if strings.ToLower(f) == "stdout" {
		file = os.Stdout
	} else {
		file, err = os.Create(f)
		if err != nil {
			return nil, err
		}
	}
	l := New(file)
	l.file = file
	return l, nil
}

// Close closes the logger.
func (l *Logger) Close() {
	close(l.logCh)
	l.waitGroup.Wait()
	if l.file != nil {
		l.file.Close()
	}
}

// Log queues a log message to be output to the logger. The message is queued only if
// the message level is less than the logger's maximum logging level.
func (l *Logger) Log(level LoggingLevel, message string) {
	if l.level < level {
		return
	}
	l.logCh <- *newLogMessage(level, message)
}

func (l *Logger) Logf(level LoggingLevel, format string, args ...any) {
	if l.level < level {
		return
	}
	l.Log(level, fmt.Sprintf(format, args...))
}

// SetMaxLevel sets the max logging level.
func (l *Logger) SetMaxLevel(level LoggingLevel) {
	l.level = Info
	l.Logf(Info, "Setting max logging level to '%s'\n", levelAsString(level))
	l.level = level
}

func levelAsString(level LoggingLevel) string {
	if level < Fatal || level > Debug {
		return fmt.Sprintf("%s:%d", LevelsAsStrings[0], level)
	}
	return LevelsAsStrings[level]
}

func newLogMessage(level LoggingLevel, msg string) *logMessage {
	return &logMessage{level: level, message: msg}
}

func (l *Logger) outputMessages() {
	for {
		msg, ok := <-l.logCh
		if !ok {
			l.waitGroup.Done()
			return
		}
		l.writer.WriteString(fmt.Sprintf("[%s]: %s", levelAsString(msg.level), msg.message))
	}
}
