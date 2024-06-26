package soapydevice

import (
	"errors"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

type sampleRateValue struct {
	rate  float64
	value string
}

var sampleRatesMap = []sampleRateValue{
	{256000.0, "256 kS/s"},
	{512000.0, "512 kS/s"},
	{1000000.0, "1.0 MS/s"},
	{1600000.0, "1.6 MS/s"},
	{2048000.0, "2.048 MS/s"},
	{2400000.0, "2.4 MS/s"},
	{2800000.0, "2.8 MS/s"},
	{3200000.0, "3.2 MS/s"},
	{4000000.0, "4.0 MS/s"},
	{5000000.0, "5.0 MS/s"},
	{6000000.0, "6.0 MS/s"},
	{7000000.0, "7.0 MS/s"},
	{8000000.0, "8.0 MS/s"},
	{9000000.0, "9.0 MS/s"},
	{10000000.0, "10.0 MS/s"},
}

// Device is the struct that holds the Soapy device properties.
type Device struct {
	sdrDevice       *device.SDRDevice
	SampleRates     []string
	sampleRate      string
	Antennas        []string
	selectedAntenna string
}

// Radio is a pointer to the SDR radio device and its properties.
var Radio *Device = nil

// Make makes a new Device object given the construction args.
//
// The device pointer is stored in a table within the Soapy API so that subsequent calls with the same arguments
// will prodcue the same device..
//
// Params:
//   - args: device key/value argument map
//
// Returns a pointer to the new Device struct or nil on error
func Make(args map[string]string) (*Device, error) {
	sdrlogger.Logf(sdrlogger.Trace, "Making device based on args: %v", args)
	dev, err := device.Make(args)
	if err == nil {
		newDevice := Device{sdrDevice: dev}
		return &newDevice, nil
	}
	sdrlogger.Logf(sdrlogger.Error, "Could not make device. Error: %v", err)
	return nil, err
}

// Unmake releases the device handle associated with the SDR device.
//
// Returns nil, or the error if the request fails.
func (dev *Device) Unmake() error {
	sdrlogger.Log(sdrlogger.Trace, "Trying to unmake device")
	err := dev.sdrDevice.Unmake()
	if err != nil {
		sdrlogger.Logf(sdrlogger.Error, "Error unmaking a Soapy Device: %v", err)
	}
	dev = nil
	return err
}

// GetSampleRateRange gets the list of available sampple rates
//
// Sets the sampleRates field of the Device to the list of available sample rates.
func (dev *Device) GetSampleRateRange() {
	sampleRates := dev.sdrDevice.GetSampleRateRange(device.DirectionRX, 0)
	sdrlogger.Logf(sdrlogger.Trace, "Sample rates = %v", sampleRates)
	var rates []string
	for _, rateRange := range sampleRates {
		minRange := rateRange.Minimum
		maxRange := rateRange.Maximum
		for _, rate := range sampleRatesMap {
			if rate.rate >= minRange && rate.rate <= maxRange {
				sdrlogger.Logf(sdrlogger.Trace, "Rate found: %v", rate.value)
				rates = append(rates, rate.value)
			}
		}
	}
	dev.SampleRates = rates
}

// GetListOfAntennas gets a list of available receive antennas for receive channel 0.
//
// Sets the antennas field of the Device to the list of available antennas.
func (dev *Device) GetListOfAntennas() {
	dev.Antennas = dev.sdrDevice.ListAntennas(device.DirectionRX, 0)
}

// SetAntenna sets the named antenna for the SDR receive channel 0.
//
// Params:
//   - antennaName: the name of the antenna to set
//
// Returns an error or nil in case of success.
func (dev *Device) SetAntenna(antennaName string) error {
	err := dev.sdrDevice.SetAntennas(device.DirectionRX, 0, antennaName)
	if err != nil {
		sdrlogger.Logf(sdrlogger.Error, "Error attempting to set antenna %v: %v", antennaName, err)
	}
	return err
}

// SetSampleRate sets the specified sample rate for SDR receive channel 0.
//
// Params:
//   - rate: the sample rate to set.
//
// Returns an error or nil on success.
func (dev *Device) SetSampleRate(rate string) error {
	var floatRate float64 = 0.0
	for _, rateValue := range sampleRatesMap {
		if rateValue.value == rate {
			floatRate = rateValue.rate
			break
		}
	}
	if floatRate == 0.0 {
		sdrlogger.Logf(sdrlogger.Error, "There is no matching sample rate for %v", rate)
		return errors.New("Programming or Logic Error: no matching sample rate for selected value")
	}
	err := dev.sdrDevice.SetSampleRate(device.DirectionRX, 0, floatRate)
	if err != nil {
		sdrlogger.Logf(sdrlogger.Error, "Error attempting to set sample rate to %v: %v", rate, err)
	}
	return err
}

// GetSamplingModeNames retrieves the sampling mode names specified by the SDR device.
//
// Returns a string array of the names, which are SDR-type specific. For example, RTL-SDR dongles, the values:
// ["Off", "I-ADC", "Q-ADC"] is returned.
func (dev *Device) GetSamplingModeNames() (string, []string) {
	settings := dev.sdrDevice.GetSettingInfo()
	sdrlogger.Logf(sdrlogger.Trace, "Settings retrieved = %v", settings)
	for _, setting := range settings {
		if setting.Key == "direct_samp" {
			sdrlogger.Logf(sdrlogger.Trace, "Sampling mode names: %v", setting.OptionNames)
			sdrlogger.Logf(sdrlogger.Trace, "Sampling mode value retrieved: %v", setting.Value)
			return setting.Value, setting.OptionNames
		}
	}
	sdrlogger.Log(sdrlogger.Error, "No sampling mode names found")
	return "", []string{}
}

// ListFrequencies lists available tunable elements in the chain.
//
// Elements are listed in order from RF to baseband.
//
// Returns a list of tunable elements by name
func (dev *Device) ListFrequencies() []string {
	return dev.sdrDevice.ListFrequencies(device.DirectionRX, 0)
}

// GetFrequencyCorrection retrieves the frequency correction.
//
// Returns the frequency correction in PPM.
func (dev *Device) GetFrequencyCorrection() float64 {
	return float64(dev.sdrDevice.GetFrequencyComponent(device.DirectionRX, 0, "CORR"))
}
