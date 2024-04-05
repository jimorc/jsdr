package soapylogging

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

// SoapyLoggingMutex helps to make sure that the logging file name does not change while writing to the file.
var SoapyLoggingMutex sync.Mutex

var soapyLogfileName string

// SoapyLoggingActive is a flag that specifies if logging should be performed.
//
// Reasons for not performing logging include:
// 1. Log file cannot be created.
// 2. You do not want to log anything.
var SoapyLoggingActive bool = false

// CreateSoapyLogfileName creates the logging file.
//
// If the file already exists, it is truncated.
// Returns error if the file cannot be created
func CreateSoapyLogfileName(name string) error {
	soapyLogfileName = name
	logFile, err := os.Create(soapyLogfileName)
	if err != nil {
	}
	err = logFile.Close()
	if err != nil {
		return err
	}
	return nil
}

// LogSoapy receives and prints Soapy messages to be logged to the log file
func LogSoapy(level sdrlogger.SDRLogLevel, message string) {
	fmt.Println(message)
	go logMessage(level, message)
}

// logMessage must be run as a goroutine
func logMessage(level sdrlogger.SDRLogLevel, message string) {
	if !SoapyLoggingActive {
		return
	}
	SoapyLoggingMutex.Lock()
	defer SoapyLoggingMutex.Unlock()

	levelStr := "Unknown"
	switch level {
	case sdrlogger.Fatal:
		levelStr = "Fatal"
	case sdrlogger.Critical:
		levelStr = "Critical"
	case sdrlogger.Error:
		levelStr = "Error"
	case sdrlogger.Warning:
		levelStr = "Warning"
	case sdrlogger.Notice:
		levelStr = "Notice"
	case sdrlogger.Info:
		levelStr = "Info"
	case sdrlogger.Debug:
		levelStr = "Debug"
	case sdrlogger.Trace:
		levelStr = "Trace"
	case sdrlogger.SSI:
		levelStr = "SSI"
	}
	logFile, err := os.OpenFile(soapyLogfileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	_, err = logFile.WriteString(fmt.Sprintf("Soapy Logged: [%v] %v\n", levelStr, message))
	if err != nil {
		log.Panic(err)
	}
}
