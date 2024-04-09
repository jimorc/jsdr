package settings

import (
	"internal/soapylogging"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

// TestSettingsUnmarshalLoggingValues tests retrieving Logging values from a JSON settings string.
func TestSettingsUnmarshalLoggingValues(t *testing.T) {
	// logging_level of 5 corresponds to sdrlogger.Notice.
	var settings string = `{
		"logging_level": 5
		}`

	testSettings := NewSettings()
	err := testSettings.Unmarshal([]byte(settings))
	if err != nil {
		t.Fatalf("TestSettingsUnmarshalLoggingValues could not unmarshal json: %v", err)
	}
	level := sdrlogger.SDRLogLevel(atomic.LoadInt64(&testSettings.LoggingLevel))
	if level != sdrlogger.Notice {
		t.Fatalf("TestSettingsUnmarshalLoggingValues could not unmarshal LoggingLevel: '%v', wanted: %v",
			soapylogging.LoggingLevelAsString(level), soapylogging.LoggingLevelAsString(sdrlogger.Notice))
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
	if atomic.LoadInt64(&testSettings.LoggingLevel) != atomic.LoadInt64(&defaultSettings.LoggingLevel) {
		t.Fatalf("TestSettingsUnmarshalEmptyJSONString has overridden Logging.LoggingLevel")
	}
}

func TestSettingsMarshal(t *testing.T) {
	settings := NewSettings()
	atomic.StoreInt64(&settings.LoggingLevel, int64(sdrlogger.Warning))

	json, err := settings.marshal()
	if err != nil {
		t.Fatal("TestSettingsMarshal could not marshal the settings struct")
	}
	settingsAsJSON := string(json)
	expected := `{"logging_level":4}`
	if !strings.HasPrefix(settingsAsJSON, expected) {
		t.Fatalf("TestSettingsMarshal did not marshal correctly: %v, expected: %v", settingsAsJSON, expected)
	}
}

func TestSettingsSaveLoad(t *testing.T) {
	settings := NewSettings()
	atomic.StoreInt64(&settings.LoggingLevel, int64(sdrlogger.Info))

	err := settings.Save()
	if err != nil {
		t.Fatalf("TestSettingsSaveLoad failed on Save call: %v", err)
	}

	newSettings := NewSettings()
	err = newSettings.Load()
	if err != nil {
		t.Fatalf("TestSettingsSaveLoad failed on Load call: %v", err)
	}

	newLoggingLevel := atomic.LoadInt64(&newSettings.LoggingLevel)
	oldLoggingLevel := atomic.LoadInt64(&settings.LoggingLevel)
	if newLoggingLevel != oldLoggingLevel {
		t.Fatalf("TestSettingsSaveLoad failed to load LoggingFile: got %v, expected %v", newLoggingLevel, oldLoggingLevel)
	}
}
