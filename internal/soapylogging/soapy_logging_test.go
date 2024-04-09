package soapylogging

import (
	"os"
	"testing"
)

func TestChangeLoggingFileName(t *testing.T) {
	fileName := "logging.log"
	badFileName := "/logging.log"
	err := ChangeLoggingFileName(fileName)
	if loggingFileName != fileName {
		t.Fatalf("Logging file name not set properly. Expected: %v, got: %v", fileName, loggingFileName)
	}
	if err != nil {
		t.Fatalf("Error creating log file '%v': %v", loggingFileName, err)
	}
	os.Remove(loggingFileName)
	err = ChangeLoggingFileName(badFileName)
	if loggingFileName != badFileName {
		t.Fatalf("Logging file name not set properly. Expected: %v, got: %v", badFileName, loggingFileName)
	}
	if err == nil {
		os.Remove(loggingFileName)
		t.Fatalf("Should not have been able to create file: %v", loggingFileName)
	}

}
