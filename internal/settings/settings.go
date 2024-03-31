package settings

import (
	"encoding/json"

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

// Unmarshal unmarshals the contents of the data byte array into the Settings struct.
func (s *Settings) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, s)
}
