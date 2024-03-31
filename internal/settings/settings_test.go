package settings

import (
	"strings"
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

func TestSettingsMarshal(t *testing.T) {
	settings := NewSettings()
	settings.Logging.LoggingFile = "log.log"
	settings.Logging.LoggingLevel = sdrlogger.Warning

	json, err := settings.marshal()
	if err != nil {
		t.Fatal("TestSettingsMarshal could not marshal the settings struct")
	}
	settingsAsJSON := string(json)
	expected := `{"logging":{"logging_file":"log.log","logging_level":4}`
	if !strings.HasPrefix(settingsAsJSON, expected) {
		t.Fatalf("TestSettingsMarshal did not marshal correctly: %v, expected: %v", settingsAsJSON, expected)
	}
}

func TestSettingsSaveLoad(t *testing.T) {
	logFile := "log.log"
	settings := NewSettings()
	settings.Logging.LoggingFile = logFile
	settings.Logging.LoggingLevel = sdrlogger.Info

	err := settings.Save()
	if err != nil {
		t.Fatalf("TestSettingsSaveLoad failed on Save call: %v", err)
	}

	newSettings := NewSettings()
	err = newSettings.Load()
	if err != nil {
		t.Fatalf("TestSettingsSaveLoad failed on Load call: %v", err)
	}

	loggingFileName := newSettings.Logging.LoggingFile
	if loggingFileName != settings.Logging.LoggingFile {
		t.Fatalf("TestSettingsSaveLoad failed to load LoggingFile: got %v, expected %v", loggingFileName, settings.Logging.LoggingFile)
	}
	loggingLevel := newSettings.Logging.LoggingLevel
	if loggingLevel != settings.Logging.LoggingLevel {
		t.Fatalf("TestSettingsSaveLoad failed to load LoggingFile: got %v, expected %v", loggingLevel, settings.Logging.LoggingLevel)
	}
}
