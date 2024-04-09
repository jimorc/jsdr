package soapylogging

import (
	"fmt"
	"log"
	"os"

	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

var loggingFileName string = os.Getenv("HOME") + "/jsdr.log"

// CreateSoapyLogFile creates the logging file.
//
// If the file already exists, it is truncated.
// Returns error if the file cannot be created
func CreateSoapyLogFile() error {
	logFile, err := os.Create(loggingFileName)
	if err != nil {
	}
	err = logFile.Close()
	if err != nil {
		return err
	}
	return nil
}

// LoggingLevelAsString converts the logging level from an int to its representative string value.
//
// If the logging level is outside its acceptable range (i.e. between Fatal and SSI), then "Unknown" is returned.
func LoggingLevelAsString(level sdrlogger.SDRLogLevel) string {
	// The level names must match the levels defined in go-soapy-sdr/pkg/sdrlogger/logger.go.
	// Since the level starts at 1 (Fatal), "Unknown" is prepended to account for value 0.
	levelStr := [10]string{"Unknown", "Fatal", "Critical", "Error", "Warning", "Notice", "Info", "Debug", "Trace", "SSI"}
	levelAsString := "Unknown"
	if level >= 0 && level <= sdrlogger.SSI {
		levelAsString = levelStr[level]
	}
	return levelAsString
}

// LogSoapy receives and prints Soapy messages to be logged to the log file
func LogSoapy(level sdrlogger.SDRLogLevel, message string) {
	logMessage(level, message)
}

// logMessage must be run as a goroutine
func logMessage(level sdrlogger.SDRLogLevel, message string) {
	levelStr := LoggingLevelAsString(level)
	logFile, err := os.OpenFile(loggingFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	_, err = logFile.WriteString(fmt.Sprintf("[%v] %v\n", levelStr, message))
	if err != nil {
		log.Panic(err)
	}
}
