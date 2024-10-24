package sdr

import (
	"fmt"
	"strings"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
)

type Antenna interface {
	GetAntennaNames(device.Direction, uint) []string
	GetCurrentAntenna(device.Direction, uint) string
}

// GetCurrentAntenna returns the currently selected RX antenna for channel 0 of the SDR.
func GetCurrentAntenna(sdrD Antenna, log *logger.Logger) string {
	antenna := sdrD.GetCurrentAntenna(device.DirectionRX, 0)
	log.Logf(logger.Debug, "Current antenna is %s\n", antenna)
	return antenna
}

// GetAntennaNames returns the list of RX antenna names for channel 0.
func GetAntennaNames(sdrD Antenna, log *logger.Logger) []string {
	antennas := sdrD.GetAntennaNames(device.DirectionRX, 0)
	var aMsg strings.Builder
	if len(antennas) == 0 {
		aMsg.WriteString("No antennas for this SDR\n")
	} else {
		aMsg.WriteString("Antennas:\n")
		for _, antenna := range antennas {
			aMsg.WriteString(fmt.Sprintf("         %s\n", antenna))
		}
		log.Log(logger.Debug, aMsg.String())
	}
	return antennas
}
