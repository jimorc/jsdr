package settings

import (
	"testing"

	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

// TestSettingsUnmarshalLoggingValues tests retrieving Logging values from a JSON settings string.
func TestSettingsUnmarshalLoggingValues(t *testing.T) {
	// logging_level of 5 corresponds to sdrlogger.Notice.
	var settings string = `{
		"logging":{
			"logging_file": "sdrLogFile.log",
			"logging_level": 5
			}
		}
	`

	testSettings := NewSettings()
	err := testSettings.Unmarshal([]byte(settings))
	if err != nil {
		t.Fatalf("TestSettingsUnmarshalLoggingValues could not unmarshal json: %v", err)
	}
	if testSettings.Logging.LoggingFile != "sdrLogFile.log" {
		t.Fatalf("TestSettTestSettingsUnmarshalLoggingValuesingsUnmarshal could not unmarshal LoggingFile: '%v', wanted: sdrLogFile.log",
			testSettings.Logging.LoggingFile)
	}
	if testSettings.Logging.LoggingLevel != sdrlogger.Notice {
		t.Fatalf("TestSettingsUnmarshalLoggingValues could not unmarshal LoggingLevel: '%v', wanted: %v",
			testSettings.Logging.LoggingLevel, sdrlogger.Notice)
	}
}

func TestSettingsUnmarshalEmptyJSONString(t *testing.T) {
	var settings string = ``
	testSettings := NewSettings()
	defaultSettings := NewSettings()
	err := testSettings.Unmarshal([]byte(settings))
	if err != nil {
		t.Fatalf("TestSettingsUnmarshalEmptyJSONString could not unmarshal an empty JSON string: %v", err)
	}
	if testSettings.Logging.LoggingFile != defaultSettings.Logging.LoggingFile {
		t.Fatalf("TestSettingsUnmarshalEmptyJSONString has overridden Logging.LoggingFile")
	}
	if testSettings.Logging.LoggingLevel != defaultSettings.Logging.LoggingLevel {
		t.Fatalf("TestSettingsUnmarshalEmptyJSONString has overridden Logging.LoggingLevel")
	}
}
