package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"internal/settings"
	"internal/soapylogging"

	"github.com/pothosware/go-soapy-sdr/pkg/device"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

func main() {
	soapylogging.SoapyLoggingActive = true
	loggingLevel := sdrlogger.Debug

	// Test log levels
	settings.JsdrSettings = settings.NewSettings()
	// Settings are different than for the jsdr app
	settings.JsdrSettings.Logging.LoggingFile = os.Getenv("HOME") + "/enumerate_sdrs.log"
	settings.JsdrSettings.Logging.LoggingLevel = loggingLevel

	soapylogging.CreateSoapyLogFile()
	sdrlogger.RegisterLogHandler(soapylogging.LogSoapy)
	sdrlogger.SetLogLevel(loggingLevel)
	sdrlogger.Log(sdrlogger.Info, "Soapy SDR")

	sdrlogger.Log(sdrlogger.Fatal, "Testing Fatal logging level")
	sdrlogger.Log(sdrlogger.Critical, "Testing Critical logging level")
	sdrlogger.Log(sdrlogger.Error, "Testing Error logging level")
	sdrlogger.Log(sdrlogger.Warning, "Testing Warning logging level")
	sdrlogger.Log(sdrlogger.Notice, "Testing Notice logging level")
	sdrlogger.Log(sdrlogger.Info, "Testing Info logging level")
	sdrlogger.Log(sdrlogger.Debug, "Testing Debug logging level")
	sdrlogger.Log(sdrlogger.Trace, "Testing Trace logging level")
	sdrlogger.Log(sdrlogger.SSI, "Testing SSI logger level")

	// List all devices
	devices := device.Enumerate(nil)
	for i, dev := range devices {
		sdrlogger.Logf(sdrlogger.Info, "Found device #%v:", i)
		for k, v := range dev {
			sdrlogger.Logf(sdrlogger.Info, "%v=%v", k, v)
		}
	}
	if len(devices) == 0 {
		sdrlogger.Logf(sdrlogger.Info, "No devices found!!")
		return
	}

	// Convert device info arguments for opening all detected devices
	deviceArgs := make([]map[string]string, len(devices))

	for i, dev := range devices {
		deviceArgs[i] = map[string]string{
			"driver": dev["driver"],
		}
	}

	// Open all devices
	devs, err := device.MakeList(deviceArgs)
	if err != nil {
		log.Panic(err)
	}
	defer func([]*device.SDRDevice) {
		// Close all devices
		err := device.UnmakeList(devs)
		if err != nil {
			log.Panic(err)
		}
		sdrlogger.Log(sdrlogger.Info, "All devices closed")
	}(devs)

	for i, dev := range devs {
		sdrlogger.Logf(sdrlogger.Info, "***************")
		sdrlogger.Logf(sdrlogger.Info, "Device: %v", devices[i]["label"])
		sdrlogger.Logf(sdrlogger.Info, "***************")

		displayDetails(dev)
		receiveSomeData(dev)
	}
}

// displayDetails displays the details and information for a device (for all its directions and channels)
func displayDetails(dev *device.SDRDevice) {
	sdrlogger.Logf(sdrlogger.Info, "Device Information")
	sdrlogger.Logf(sdrlogger.Info, "***************")

	// Print hardware info for the device
	displayHardwareInfo(dev)

	// GPIO
	displayGPIOBanks(dev)

	// Settings
	displaySettingInfo(dev)

	// UARTs
	displayUARTs(dev)

	// Clocking
	displayMasterClockRate(dev)
	displayClockSources(dev)

	// Registers
	displayRegisters(dev)

	// Device Sensor
	displaySensors(dev)

	// Time Sources
	displayTimeSources(dev)

	// Direction details
	displayDirectionDetails(dev, device.DirectionTX)
	displayDirectionDetails(dev, device.DirectionRX)
}

// displayHardwareInfo prints hardware info for the specified device
func displayHardwareInfo(dev *device.SDRDevice) {
	sdrlogger.Logf(sdrlogger.Info, "DriverKey: %v", dev.GetDriverKey())
	sdrlogger.Logf(sdrlogger.Info, "HardwareKey: %v", dev.GetHardwareKey())
	hardwareInfo := dev.GetHardwareInfo()
	if len(hardwareInfo) > 0 {
		for k, v := range hardwareInfo {
			sdrlogger.Logf(sdrlogger.Info, "HardwareInfo: %v: %v", k, v)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "HardwareInfo: [none]")
	}
}

// displayGPIOBanks prints GPIO bank info for the specified device
func displayGPIOBanks(dev *device.SDRDevice) {
	banks := dev.ListGPIOBanks()
	if len(banks) > 0 {
		for i, bank := range banks {
			sdrlogger.Logf(sdrlogger.Info, "GPIO Bank#%d: %v", i, bank)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "GPIO Banks: [none]")
	}
}

// displaySettingInfo prints a device's setting information
func displaySettingInfo(dev *device.SDRDevice) {
	settings := dev.GetSettingInfo()
	if len(settings) > 0 {
		for i, setting := range settings {
			sdrlogger.Logf(sdrlogger.Info, "Setting#%d:", i)
			displaySettingValues(setting)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Settings: [none]")
	}
}

// displaySettingValues prints each setting value
func displaySettingValues(setting device.SDRArgInfo) {
	sdrlogger.Logf(sdrlogger.Info, "  key: %v", setting.Key)
	sdrlogger.Logf(sdrlogger.Info, "  value: %v", setting.Value)
	sdrlogger.Logf(sdrlogger.Info, "  name: %v", setting.Name)
	sdrlogger.Logf(sdrlogger.Info, "  description: %v", setting.Description)
	sdrlogger.Logf(sdrlogger.Info, "  unit: %v", setting.Unit)
	var argType string = "unknown type"
	switch setting.Type {
	case device.ArgInfoBool:
		argType = "bool"
	case device.ArgInfoInt:
		argType = "integer"
	case device.ArgInfoFloat:
		argType = "float"
	case device.ArgInfoString:
		argType = "string"
	}
	sdrlogger.Logf(sdrlogger.Info, "  type: %v", argType)
	sdrlogger.Logf(sdrlogger.Info, "  range: %v", setting.Range.ToString())
	numOptions := setting.NumOptions
	if numOptions > 0 {
		sdrlogger.Logf(sdrlogger.Info, "  options: %v", setting.Options)
		sdrlogger.Logf(sdrlogger.Info, "  option names: %v", setting.OptionNames)
	} else {
		sdrlogger.Logf(sdrlogger.Info, "  options: [none]")
		sdrlogger.Logf(sdrlogger.Info, "  option names: [none]")
	}
}

// displayUARTs prints a devices's UARTs
func displayUARTs(dev *device.SDRDevice) {
	uarts := dev.ListUARTs()
	if len(uarts) > 0 {
		for i, uart := range uarts {
			sdrlogger.Logf(sdrlogger.Info, "UARTs#%d:%v", i, uart)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "UARTs: [none]")
	}
}

// displayMasterClockRate prints a device's master clock rate and clock ranges
func displayMasterClockRate(dev *device.SDRDevice) {
	sdrlogger.Logf(sdrlogger.Info, "Master Clock Rate: %v", dev.GetMasterClockRate())
	clockRanges := dev.GetMasterClockRates()
	if len(clockRanges) > 0 {
		sdrlogger.Logf(sdrlogger.Info, "Master Clock Rate Ranges:")
		for i, clockRange := range clockRanges {
			sdrlogger.Logf(sdrlogger.Info, "  Range#%d: %v", i, clockRange)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Clock Rate Ranges: [none]")
	}
}

// displayClockSources prints a device's clock sources
func displayClockSources(dev *device.SDRDevice) {
	clockSources := dev.ListClockSources()
	if len(clockSources) > 0 {
		sdrlogger.Logf(sdrlogger.Info, "Clock Sources:")
		for i, clockSource := range clockSources {
			sdrlogger.Logf(sdrlogger.Info, "  Clock Source#%d: %v", i, clockSource)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Clock Sources: [none]")
	}
}

// displayRegisters prints a device's registers
func displayRegisters(dev *device.SDRDevice) {
	registers := dev.ListRegisterInterfaces()
	if len(registers) > 0 {
		sdrlogger.Logf(sdrlogger.Info, "Registers:")
		for i, register := range registers {
			sdrlogger.Logf(sdrlogger.Info, "  Register#%d: %v", i, register)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Registers: [none]")
	}
}

// displaySensors prints a device's sensors
func displaySensors(dev *device.SDRDevice) {
	sensors := dev.ListSensors()
	if len(sensors) > 0 {
		sdrlogger.Logf(sdrlogger.Info, "Sensors:")
		for i, sensor := range sensors {
			sdrlogger.Logf(sdrlogger.Info, "  Sensor#%d: %v", i, sensor)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Sensors: [none]")
	}
}

// displayTimeSources lists all of a device's time sources and hardware time if any
func displayTimeSources(dev *device.SDRDevice) {
	timeSources := dev.ListTimeSources()
	if len(timeSources) > 0 {
		sdrlogger.Logf(sdrlogger.Info, "Time Sources:")
		for i, timeSource := range timeSources {
			sdrlogger.Logf(sdrlogger.Info, "  Time Source#%d: %v", i, timeSource)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Time Sources: [none]")
	}

	hasHardwareTime := dev.HasHardwareTime("")
	sdrlogger.Logf(sdrlogger.Info, "Has Hardware Time: %v", hasHardwareTime)
	if hasHardwareTime {
		sdrlogger.Logf(sdrlogger.Info, "  Hardware Time: %v", dev.GetHardwareTime(""))
	}
}

// displayDirectionDetails prints info about TX and RX channels
func displayDirectionDetails(dev *device.SDRDevice, direction device.Direction) {
	if direction == device.DirectionTX {
		sdrlogger.Logf(sdrlogger.Info, "Direction TX")
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Direction RX")
	}

	frontEndMapping := dev.GetFrontendMapping(direction)
	if len(frontEndMapping) > 0 {
		sdrlogger.Logf(sdrlogger.Info, "  FrontendMapping: %v", frontEndMapping)
	} else {
		sdrlogger.Logf(sdrlogger.Info, "  FrontendMapping: [none]")
	}

	numChannels := dev.GetNumChannels(direction)
	sdrlogger.Logf(sdrlogger.Info, "  Number of Channels: %v", numChannels)

	for channel := uint(0); channel < numChannels; channel++ {
		displayDirectionChannelDetails(dev, direction, channel)
	}
}

// displayDirectionChannelDetails prints out details and info of a device / direction / channel
func displayDirectionChannelDetails(dev *device.SDRDevice, direction device.Direction, channel uint) {
	// Settings
	settings := dev.GetChannelSettingInfo(direction, channel)
	if len(settings) > 0 {
		for i, setting := range settings {
			sdrlogger.Logf(sdrlogger.Info, "Channel#%d Setting#%d Banks: %v", channel, i, setting)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Settings: [none]", channel)
	}

	// Channel
	channelInfo := dev.GetChannelInfo(direction, channel)
	if len(channelInfo) > 0 {
		for k, v := range channelInfo {
			sdrlogger.Logf(sdrlogger.Info, "Channel#%d ChannelInfo: {%v: %v}", channel, k, v)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d ChannelInfo: [none]", channel)
	}

	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Fullduplex: %v", channel, dev.GetFullDuplex(direction, channel))

	// Antenna
	antennas := dev.ListAntennas(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d NumAntennas: %v", channel, len(antennas))

	for i, antenna := range antennas {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Antenna#%d: %v", channel, i, antenna)
	}

	// Bandwidth
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Baseband filter width: %v Hz", channel, dev.GetBandwidth(direction, channel))

	bandwidthRanges := dev.GetBandwidthRanges(direction, channel)
	for i, bandwidthRange := range bandwidthRanges {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Baseband filter#%d: %v", channel, i, bandwidthRange)
	}

	// Gain
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d HasGainMode (Automatic gain possible): %v", channel, dev.HasGainMode(direction, channel))
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d GainMode (Automatic gain enabled): %v", channel, dev.GetGainMode(direction, channel))
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Gain: %v", channel, dev.GetGain(direction, channel))
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d GainRange: %v", channel, dev.GetGainRange(direction, channel))
	gainElements := dev.ListGains(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d NumGainElements: %v", channel, len(gainElements))

	for i, gainElement := range gainElements {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Gain Element#%d Name: %v", channel, i, gainElement)
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Gain Element#%d Value: %v", channel, i, dev.GetGainElement(direction, channel, gainElement))
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Gain Element#%d Range: %v", channel, i, dev.GetGainElementRange(direction, channel, gainElement).ToString())
	}

	// SampleRate
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Sample Rate: %v", channel, dev.GetSampleRate(direction, channel))
	sampleRateRanges := dev.GetSampleRateRange(direction, channel)
	for i, sampleRateRange := range sampleRateRanges {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Sample Rate Range#%d: %v", channel, i, sampleRateRange.ToString())
	}

	// Frequencies
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Frequency: %v", channel, dev.GetFrequency(direction, channel))
	frequencyRanges := dev.GetFrequencyRange(direction, channel)
	for i, frequencyRange := range frequencyRanges {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Frequency Range#%d: %v", channel, i, frequencyRange.ToString())
	}

	frequencyArgsInfos := dev.GetFrequencyArgsInfo(direction, channel)

	if len(frequencyArgsInfos) > 0 {
		for i, argInfo := range frequencyArgsInfos {
			sdrlogger.Logf(sdrlogger.Info, "Channel#%d Frequency ArgInfo#%d: %v", channel, i, argInfo.ToString())
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Frequency ArgInfo: [none]", channel)
	}

	frequencyComponents := dev.ListFrequencies(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d NumFrequencyComponents: %v", channel, len(frequencyComponents))

	for i, frequencyComponent := range frequencyComponents {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Frequency Component#%d Name: %v", channel, i, frequencyComponent)
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Frequency Component#%d Frequency: %v", channel, i,
			dev.GetFrequencyComponent(direction, channel, frequencyComponent))
	}

	// Stream
	streamFormats := dev.GetStreamFormats(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Formats: %v", channel, streamFormats)
	nativeStreamFormat, nativeStreamFullScale := dev.GetNativeStreamFormat(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Native Format: %v", channel, nativeStreamFormat)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Native FullScale: %v", channel, nativeStreamFullScale)

	streamArgsInfos := dev.GetStreamArgsInfo(direction, channel)
	if len(streamArgsInfos) > 0 {
		for i, argInfo := range streamArgsInfos {
			sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream ArgInfo#%d: %v", channel, i, argInfo.ToString())
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream ArgInfo: [none]", channel)
	}

	// Frontend correctiion
	available := dev.HasDCOffsetMode(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Correction Auto DC correction available: %v", channel, available)
	if available {
		offsetMode := dev.GetDCOffsetMode(direction, channel)
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Correction Auto DEC correction: %v", channel, offsetMode)
	}

	available = dev.HasDCOffset(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Correction DC Correction available: %v", channel, available)
	if available {
		I, Q, err := dev.GetDCOffset(direction, channel)
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Correction DC correction I: %v, Q: %v, err: %v", channel, I, Q, err)
	}

	available = dev.HasIQBalance(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Correction IQ Balance available: %v", channel, available)
	if available {
		I, Q, err := dev.GetIQBalance(direction, channel)
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Correction IQ Balnance I: %v, Q: %v, err: %v", channel, I, Q, err)
	}

	available = dev.HasFrequencyCorrection(direction, channel)
	sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Correction Frequency correction available: %v", channel, available)
	if available {
		frequencyCorrection := dev.GetFrequencyCorrection(direction, channel)
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Stream Correction Frequency correction: %v PPM", channel, frequencyCorrection)
	}

	// Channel Sensor
	sensors := dev.ListChannelSensors(direction, channel)
	if len(sensors) > 0 {
		for i, sensor := range sensors {
			sdrlogger.Logf(sdrlogger.Info, "Channel#%d Sensor#%d: %v", channel, i, sensor)
		}
	} else {
		sdrlogger.Logf(sdrlogger.Info, "Channel#%d Sensors: [none]", channel)
	}
}

// receiveSomeData receives CS8 data from stream 0
func receiveSomeData(dev *device.SDRDevice) {
	sdrlogger.Logf(sdrlogger.Info, "---------------")
	sdrlogger.Logf(sdrlogger.Info, "Data Reception")
	sdrlogger.Logf(sdrlogger.Info, "---------------")

	// Apply settings
	if err := dev.SetSampleRate(device.DirectionRX, 0, 1e6); err != nil {
		log.Fatal(fmt.Printf("SetSampleRate fail: error: %v\n", err))
	}
	if err := dev.SetFrequency(device.DirectionRX, 0, 99.9e6, nil); err != nil {
		log.Fatal(fmt.Printf("SetFrequency fail: error: %v\n", err))
	}

	stream, err := dev.SetupSDRStreamCS8(device.DirectionRX, []uint{0}, nil)

	if err != nil {
		log.Fatal(fmt.Printf("SetupStream fail: error: %v\n", err))
	}

	if err := stream.Activate(0, 0, 0); err != nil {
		log.Fatal(fmt.Printf("Stream Activate fail: error: %v\n", err))
	}

	mtu := stream.GetMTU()
	sdrlogger.Logf(sdrlogger.Info, "Stream MTU: %v", mtu)
	numBuffers := stream.GetNumDirectAccessBuffers()
	sdrlogger.Logf(sdrlogger.Info, "NumDirectAccessBuffers: %v", numBuffers)

	buffers := make([][]int8, 1)
	buffers[0] = make([]int8, 1024)
	flags := make([]int, 1)

	for i := 0; i < 10; i++ {
		timeNs, numElementsRead, err := stream.Read(buffers, 511, flags, 1000000)
		var flag int = flags[0]
		streamFlags := buildStreamFlagsString(flag)
		sdrlogger.Logf(sdrlogger.Info, "flags=%v", streamFlags)
		sdrlogger.Logf(sdrlogger.Info, "numElemsRead=%v, timeNS=%v, err=%v", numElementsRead, timeNs, err)
		if err == nil {
			for j := uint(0); j < numElementsRead; j += 8 {
				sdrlogger.Log(sdrlogger.Info, buildDataLine(numElementsRead, 8, j, buffers[0]))
			}
		}
	}

	if err := stream.Deactivate(0, 1000000); err != nil {
		log.Fatal(fmt.Printf("Stream Deactivate fail: error: %v\n", err))
	}

	if err := stream.Close(); err != nil {
		log.Fatal(fmt.Printf("Stream close fail: error: %v\n", err))
	}
}

// buildStreamFlagsString builds a string of stream flag names.
func buildStreamFlagsString(flag int) string {
	var haveFlag = false
	var flagStringBuilder strings.Builder
	addFlagStringToStringBuilder(flag, "EndBurst", device.StreamFlagEndBurst, &haveFlag, &flagStringBuilder)
	addFlagStringToStringBuilder(flag, "HasTime", device.StreamFlagHasTime, &haveFlag, &flagStringBuilder)
	addFlagStringToStringBuilder(flag, "EndAbrupt", device.StreamFlagEndAbrupt, &haveFlag, &flagStringBuilder)
	addFlagStringToStringBuilder(flag, "OnePacket", device.StreamFlagOnePacket, &haveFlag, &flagStringBuilder)
	addFlagStringToStringBuilder(flag, "MoreFragments", device.StreamFlagMoreFragments, &haveFlag, &flagStringBuilder)
	addFlagStringToStringBuilder(flag, "WaitTrigger", device.StreamFlagWaitTrigger, &haveFlag, &flagStringBuilder)
	if flagStringBuilder.Len() > 0 {
		return flagStringBuilder.String()
	} else {
		return "[none]"
	}
}

func buildDataLine(bufferSize uint, lineSize uint, startElement uint, data []int8) string {
	var dataBuilder strings.Builder
	elementsToPrint := lineSize
	if startElement+lineSize > bufferSize {
		elementsToPrint = bufferSize - startElement
	}
	for element := startElement; element < startElement+elementsToPrint; element++ {
		dataBuilder.WriteString(fmt.Sprintf("{%v, %v} ", data[2*element], data[2*element+1]))
	}
	return dataBuilder.String()
}

// addFlagStringToStringBuilder adds stream flag name to string.Builder if flag is set.
func addFlagStringToStringBuilder(flag int, flagName string, testFlag device.StreamFlag, haveFlags *bool, flagString *strings.Builder) {
	if flag&int(testFlag) == int(testFlag) {
		if *haveFlags {
			flagString.WriteString(", ")
		}
		*haveFlags = true
		flagString.WriteString(flagName)
	}
}
