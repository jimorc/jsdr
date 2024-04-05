package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

// JsdrSettings is a global variable that holds program settings.
//
// The first thing that the program should do is initialize this variable by calling settings.NewSettings().
var JsdrSettings *Settings

// SettingsMutex is used to control access to the JsdrSettings struct
var SettingsMutex sync.Mutex

// LoggingSettings contains SDRLogging related settings.
type LoggingSettings struct {
	LoggingFile  string                `json:"logging_file,omitempty"`
	LoggingLevel sdrlogger.SDRLogLevel `json:"logging_level,omitempty"`
}

// Settings contains values that are shared between executions of go_sdr.
type Settings struct {
	Logging LoggingSettings `json:"logging,omitempty"`
}

// NewSettings creates a new default Settings struct.
func NewSettings() *Settings {
	var settings Settings
	settings.Logging.LoggingFile = os.Getenv("HOME") + "/jsdr.log"
	settings.Logging.LoggingLevel = sdrlogger.Debug

	return &settings
}

// Load opens the JSON formatted file *filename* and unmarshals it into the Settings struct.
func (s *Settings) Load() error {
	file, err := os.Open("JsdrSettings.Logging.LoggingFile")
	if err != nil {
		return err
	}
	defer file.Close()

	fInfo, err := os.Stat(JsdrSettings.Logging.LoggingFile)
	if err != nil {
		return err
	}
	fileSize := fInfo.Size()
	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		return err
	}
	s.Unmarshal(data)
	return nil
}

// Save writes the JSON-formatted settings to *filename*.
func (s *Settings) Save() error {
	file, err := os.Create(JsdrSettings.Logging.LoggingFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := s.marshal()
	if err != nil {
		return err
	}
	numBytesWritten, err := file.Write(data)
	if err != nil {
		return err
	}
	if numBytesWritten != len(data) {
		return fmt.Errorf("Settings.Save wrote %d of %d bytes", numBytesWritten, len(data))
	}
	return nil
}

// Unmarshal unmarshals the contents of the data byte array into the Settings struct.
func (s *Settings) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, s)
}

// marshal marshals the contents of the Settings struct.
func (s *Settings) marshal() ([]byte, error) {
	return json.Marshal(s)
}
