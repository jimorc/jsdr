package settings

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

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
	settings.Logging.LoggingFile = "go_sdr.log"
	settings.Logging.LoggingLevel = sdrlogger.Info

	return &settings
}

// Load opens the JSON formatted file *filename* and unmarshals it into the Settings struct.
func (s *Settings) Load() error {
	settingsFileName := os.Getenv("HOME") + "/.jsdr"
	file, err := os.Open(settingsFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fInfo, err := os.Stat(settingsFileName)
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
	settingsFileName := os.Getenv("HOME") + "/.jsdr"
	file, err := os.Create(settingsFileName)
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
