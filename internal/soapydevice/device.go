package soapydevice

import (
	"github.com/pothosware/go-soapy-sdr/pkg/device"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

// Device is the struct that holds the Soapy device properties.
type Device struct {
	sdrDevice *device.SDRDevice
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
