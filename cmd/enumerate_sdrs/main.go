package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jimorc/jsdr/internal/logger"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	logLevel, logFile := parseCommandLine()

	log := initLogfile(logLevel, logFile)
	defer log.Close()

	devices := device.Enumerate(nil)
	log.Logf(logger.Info, "Found %d attached SDR(s)\n", len(devices))

	for i, dev := range devices {
		var devInfo strings.Builder
		devInfo.WriteString(fmt.Sprintf("Device %d\n", i))
		for k, v := range dev {
			devInfo.WriteString(fmt.Sprintf("         %s: %s\n", k, v))
		}
		log.Log(logger.Info, devInfo.String())

		// Open device
		log.Logf(logger.Debug, "Making device with label: '%s'\n", dev["label"])
		sdr, err := device.Make(dev)
		if err != nil {
			log.Logf(logger.Error,
				"Unable to make device with label: %s: %s\n", dev["label"], err.Error())
		}
		if sdr == nil {
			log.Log(logger.Error, "Could not make SDR\n")
		}
		defer func() {
			err := sdr.Unmake()
			if err != nil {
				log.Logf(logger.Error,
					"Could not Unmake SDR with label: %s: %s\n",
					dev["label"], err.Error())
				fmt.Println("Unable to Unmake a device. See log file for more info.")
				os.Exit(1)
			}
			log.Logf(logger.Debug,
				"Device with label: `%s` was unmade.\n",
				dev["label"])
		}()

		logHardwareInfo(sdr, log)
		logGPIOBanks(sdr, log)
		logSettingInfo(sdr, log)
		logUARTs(sdr, log)
		logMasterClockRate(sdr, log)
		logClockSources(sdr, log)
		logRegisters(sdr, log)
		logSensors(sdr, log)
		logTimeSources(sdr, log)
		logDirectionDetails(sdr, device.DirectionTX, log)
		logDirectionDetails(sdr, device.DirectionRX, log)
	}
}

func parseCommandLine() (logger.LoggingLevel, string) {
	pflag.Bool("debug", false, "Log debug information")
	pflag.String("out", os.Getenv("HOME")+"/enumerate_sdrs.log", "Log filename. If 'stdout', messages are logged to 'stdout.")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	debug := viper.GetBool("debug")
	logFile := viper.GetString("out")
	logLevel := logger.Info
	if debug {
		logLevel = logger.Debug
	}
	return logLevel, logFile
}

func initLogfile(level logger.LoggingLevel, fileName string) *logger.Logger {
	log, err := logger.NewFileLogger(fileName)
	if err != nil {
		fmt.Printf("Error trying to open log file '%s': %s\n", fileName, err.Error())
		os.Exit(1)
	}
	log.SetMaxLevel(level)
	return log
}

func logHardwareInfo(sdr *device.SDRDevice, log *logger.Logger) {
	var hwInfo strings.Builder
	hwInfo.WriteString(fmt.Sprintln("Hardware Info:"))
	hwInfo.WriteString(fmt.Sprintf("         Driver Key: %s\n", sdr.GetDriverKey()))
	hwInfo.WriteString(fmt.Sprintf("         Hardware Key: %s\n", sdr.GetHardwareKey()))
	hardwareInfo := sdr.GetHardwareInfo()

	for k, v := range hardwareInfo {
		hwInfo.WriteString(fmt.Sprintf("         %s: %s\n", k, v))
	}
	log.Log(logger.Info, hwInfo.String())
}

func logGPIOBanks(sdr *device.SDRDevice, log *logger.Logger) {
	banks := sdr.ListGPIOBanks()
	if len(banks) == 0 {
		log.Log(logger.Info, "GPIO Banks: none\n")
	} else {
		var gpioBanks strings.Builder
		gpioBanks.WriteString("GPIO Banks:\n")
		for i, bank := range banks {
			gpioBanks.WriteString(fmt.Sprintf("         GPIO Bank#%d: %v\n", i, bank))
		}
		log.Log(logger.Info, gpioBanks.String())
	}
}

func logSettingInfo(sdr *device.SDRDevice, log *logger.Logger) {
	SDRSettings := sdr.GetSettingInfo()
	if len(SDRSettings) == 0 {
		log.Log(logger.Info, "Settings: none")
	} else {
		var settings strings.Builder
		for i, set := range SDRSettings {
			if i == 0 {
				settings.WriteString(fmt.Sprintf("Setting%d:\n", i))
			} else {
				settings.WriteString(fmt.Sprintf("        Setting%d:\n", i))
			}
			settings.WriteString(fmt.Sprintf("         key: %s\n", set.Key))
			settings.WriteString(fmt.Sprintf("         value: %s\n", set.Value))
			settings.WriteString(fmt.Sprintf("         name: %s\n", set.Name))
			settings.WriteString(fmt.Sprintf("         description: %s\n", set.Description))
			settings.WriteString(fmt.Sprintf("         unit: %s\n", set.Unit))
			argType := "unknown type"
			switch set.Type {
			case device.ArgInfoBool:
				argType = "bool"
			case device.ArgInfoInt:
				argType = "int"
			case device.ArgInfoFloat:
				argType = "float"
			case device.ArgInfoString:
				argType = "string"
			}
			settings.WriteString(fmt.Sprintf("         type: %s\n", argType))
			settings.WriteString(fmt.Sprintf("         range: %v\n", set.Range.ToString()))
			numOptions := set.NumOptions
			if numOptions == 0 {
				settings.WriteString(fmt.Sprintln("         options: none"))
				settings.WriteString(fmt.Sprintln("         option names: none"))
			} else {
				settings.WriteString(fmt.Sprintln("         options:"))
				for _, opt := range set.Options {
					settings.WriteString(fmt.Sprintf("            %s\n", opt))
				}
				settings.WriteString(fmt.Sprintln("            option names:"))
				for _, name := range set.OptionNames {
					settings.WriteString(fmt.Sprintf("             %s\n", name))
				}
			}
		}
		log.Log(logger.Info, settings.String())
	}
}

func logUARTs(sdr *device.SDRDevice, log *logger.Logger) {
	uarts := sdr.ListUARTs()
	if len(uarts) == 0 {
		log.Log(logger.Info, "UARTs: none\n")
	} else {
		var umsg strings.Builder
		umsg.WriteString("UARTs:\n")
		for i, uart := range uarts {
			umsg.WriteString(fmt.Sprintf("         UART#%d: %s", i, uart))
		}
		log.Log(logger.Info, umsg.String())
	}
}

func logMasterClockRate(sdr *device.SDRDevice, log *logger.Logger) {
	clockRates := sdr.GetMasterClockRates()
	if len(clockRates) == 0 {
		log.Log(logger.Info, "Master Clock Rates: none\n")
	} else {
		log.Logf(logger.Info, "Master Clock Rate: %f\n", sdr.GetMasterClockRate())
		var rMsg strings.Builder
		rMsg.WriteString("Master Clock Rates:\n")
		for _, rate := range clockRates {
			rMsg.WriteString(fmt.Sprintf("         %v\n", rate))
		}
		log.Log(logger.Info, rMsg.String())
	}
}

func logClockSources(sdr *device.SDRDevice, log *logger.Logger) {
	sources := sdr.ListClockSources()
	if len(sources) == 0 {
		log.Log(logger.Info, "Clock Sources: none\n")
	} else {
		var sMsg strings.Builder
		sMsg.WriteString("Clock Sources:\n")
		for i, source := range sources {
			sMsg.WriteString(fmt.Sprintf("         Source#%d: %s\n", i, source))
		}
		log.Log(logger.Info, sMsg.String())
	}
}

func logRegisters(sdr *device.SDRDevice, log *logger.Logger) {
	registers := sdr.ListRegisterInterfaces()
	if len(registers) == 0 {
		log.Log(logger.Info, "Registers: none\n")
	} else {
		var rMsg strings.Builder
		rMsg.WriteString("Registers:\n")
		for i, register := range registers {
			rMsg.WriteString(fmt.Sprintf("         Register#%d: %s\n", i, register))
		}
		log.Log(logger.Info, rMsg.String())
	}
}

func logSensors(sdr *device.SDRDevice, log *logger.Logger) {
	sensors := sdr.ListSensors()
	if len(sensors) == 0 {
		log.Log(logger.Info, "Sensors: none\n")
	} else {
		var sMsg strings.Builder
		sMsg.WriteString("Sensors:\n")
		for i, sensor := range sensors {
			sMsg.WriteString(fmt.Sprintf("         Sensor#%d: %s\n", i, sensor))
		}
		log.Log(logger.Info, sMsg.String())
	}
}

func logTimeSources(sdr *device.SDRDevice, log *logger.Logger) {
	sources := sdr.ListTimeSources()
	if len(sources) == 0 {
		log.Log(logger.Info, "Time Sources: none\n")
	} else {
		var tMsg strings.Builder
		tMsg.WriteString("Time Sources:\n")
		for i, source := range sources {
			tMsg.WriteString(fmt.Sprintf("         Time Source#%d: %s\n", i, source))
		}
		log.Log(logger.Info, tMsg.String())
	}

	hasHardwareTime := sdr.HasHardwareTime("")
	log.Logf(logger.Info, "Has Hardware Time: %v\n", hasHardwareTime)
	if hasHardwareTime {
		log.Logf(logger.Info, "Hardware Time: %d ns\n", sdr.GetHardwareTime(""))
		curTime := time.Now().UTC().Nanosecond()
		log.Logf(logger.Info, "Setting Hardware Time to %d\n", curTime)
		sdr.SetHardwareTime(uint(curTime), "")
		log.Logf(logger.Info, "Hardware Time Now: %d\n", sdr.GetHardwareTime(""))
		log.Log(logger.Info, "Waiting 1 second\n")
		time.Sleep(time.Second)
		log.Logf(logger.Info, "Hardware Time Now: %d\n", sdr.GetHardwareTime(""))
	}
}

func logDirectionDetails(sdr *device.SDRDevice, direction device.Direction, log *logger.Logger) {
	if direction == device.DirectionTX {
		log.Log(logger.Info, "Direction TX\n")
	} else {
		log.Log(logger.Info, "Direction RX\n")
	}

	frontendMapping := sdr.GetFrontendMapping(direction)
	if len(frontendMapping) == 0 {
		log.Log(logger.Info, "Frontend Mapping: none\n")
	} else {
		log.Logf(logger.Info, "Frontend Mapping: %s\n", frontendMapping)
	}

	numChannels := sdr.GetNumChannels(direction)
	log.Logf(logger.Info, "Number of channels: %d\n", numChannels)
	for ch := uint(0); ch < numChannels; ch++ {
		logDirectionChannelDetails(sdr, direction, ch, log)
	}
}

func logDirectionChannelDetails(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	logChannelSettingsInfo(sdr, direction, channel, log)
	logChannelInfo(sdr, direction, channel, log)
	exerciseAntennas(sdr, direction, channel, log)
	exerciseChannelBandwidth(sdr, direction, channel, log)
	exerciseGain(sdr, direction, channel, log)
	exerciseSampleRate(sdr, direction, channel, log)
	exerciseFrequencies(sdr, direction, channel, log)
	logStreamFormatsAndInfo(sdr, direction, channel, log)
	exerciseFrontend(sdr, direction, channel, log)
	exerciseChannelSensors(sdr, direction, channel, log)
	exerciseStream0(sdr, direction, channel, log)
}

func logChannelSettingsInfo(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	settings := sdr.GetChannelSettingInfo(direction, channel)
	if len(settings) == 0 {
		log.Logf(logger.Info, "Channel#%d Settings: none\n", channel)
	} else {
		var sMsg strings.Builder
		sMsg.WriteString(fmt.Sprintf("Channel#%d Settings:\n", channel))
		for i, setting := range settings {
			sMsg.WriteString(fmt.Sprintf("         Channel#%d Setting#%d: %v\n", channel, i, setting))
		}
		log.Log(logger.Info, sMsg.String())
	}
}

func logChannelInfo(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	channelInfo := sdr.GetChannelInfo(direction, channel)
	if len(channelInfo) == 0 {
		log.Logf(logger.Info, "Channel#%d ChannelInfo: none\n", channel)
	} else {
		var infoMsg strings.Builder
		infoMsg.WriteString(fmt.Sprintf("Channel#%d ChannelInfo:\n", channel))
		for k, v := range channelInfo {
			infoMsg.WriteString(fmt.Sprintf("         %s: %s\n", k, v))
		}
	}
}

func exerciseAntennas(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	antennas := sdr.ListAntennas(direction, channel)
	if len(antennas) == 0 {
		log.Log(logger.Info, "Antennas: none\n")
	} else {
		log.Logf(logger.Info, "Selected antenna: %s\n",
			sdr.GetAntennas(direction, channel))
		var aMsg strings.Builder
		aMsg.WriteString("Antennas:\n")
		for i, antenna := range antennas {
			aMsg.WriteString(fmt.Sprintf("         Antenna#%d: %s\n", i, antenna))
			aMsg.WriteString(fmt.Sprintf("         Setting antenna to %s\n", antenna))
			sdr.SetAntennas(direction, channel, antenna)
			aMsg.WriteString(fmt.Sprintf("         Selected antenna now %s\n", sdr.GetAntennas(direction, channel)))
		}
		log.Log(logger.Info, aMsg.String())
	}
}

func exerciseChannelBandwidth(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	log.Logf(logger.Info,
		"Channel#%d Baseband filter width: %.0f Hz\n", channel, sdr.GetBandwidth(direction, channel))

	bandwidthRanges := sdr.GetBandwidthRanges(direction, channel)
	if len(bandwidthRanges) == 0 {
		log.Logf(logger.Info, "Channel#%d Bandwidth Ranges: none\n", channel)
	} else {
		var bMsg strings.Builder
		bMsg.WriteString(fmt.Sprintf("Channel#%d Bandwidth Ranges:\n", channel))
		for i, bRange := range bandwidthRanges {
			bMsg.WriteString(fmt.Sprintf("         Bandwidth Range#%d: %v\n", i, bRange))
		}
		log.Log(logger.Info, bMsg.String())

		log.Log(logger.Info, "Setting bandwidth to one half first range\n")
		err := sdr.SetBandwidth(direction, channel, bandwidthRanges[0].Maximum/2.0)
		if err != nil {
			log.Logf(logger.Error, "Error encountered while trying to set bandwidth: %sn",
				err.Error())
		}
		log.Logf(logger.Info, "Bandwidth is now %.0f\n", sdr.GetBandwidth(direction, channel))
	}
}

func exerciseGain(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	hasAutoGainMode := sdr.HasGainMode(direction, channel)
	log.Logf(logger.Info, "Channel#%d HasGainMode (Automatic gain possible): %v\n",
		channel, hasAutoGainMode)
	if hasAutoGainMode {
		autoGainEnabled := sdr.GetGainMode(direction, channel)
		log.Logf(logger.Info, "Channel#%d Automatic Gain Enabled: %v\n",
			channel, autoGainEnabled)
		log.Logf(logger.Info, "Current gain = %f\n",
			sdr.GetGain(direction, channel))
		log.Log(logger.Info, "Toggling auto gain\n")
		err := sdr.SetGainMode(direction, channel, !autoGainEnabled)
		if err != nil {
			log.Logf(logger.Error, "Error in call to SetGainMode: %s\n", err.Error())
			return
		}
		log.Logf(logger.Info, "Channel#%d Automatic Gain Enabled now: %v\n",
			channel, sdr.GetGainMode(direction, channel))
		log.Logf(logger.Info, "Current gain = %f\n",
			sdr.GetGain(direction, channel))
	}
	gains := sdr.ListGains(direction, channel)
	log.Logf(logger.Info, "Number of gain elements: %d\n", len(gains))
	if len(gains) > 0 {
		var gMsg strings.Builder
		gMsg.WriteString(fmt.Sprintf("Channel#%d Gain Elements:\n", channel))
		for _, gain := range gains {
			gMsg.WriteString(fmt.Sprintf("         Element: %s\n", gain))
			gMsg.WriteString(fmt.Sprintf("             Range: %v\n", sdr.GetGainElementRange(direction, channel, gain)))
		}
		log.Log(logger.Info, gMsg.String())

		err := sdr.SetGainMode(direction, channel, false)
		if err != nil {
			log.Logf(logger.Error, "Error setting auto gain off: %s\n", err.Error())
			return
		} else {
			log.Log(logger.Info, "Have set auto gain off\n")
		}
		for _, gain := range gains {
			log.Logf(logger.Info, "Setting gain for element: %s to 20 db\n", gain)
			err := sdr.SetGainElement(direction, channel, gain, 20.0)
			if err != nil {
				log.Logf(logger.Error, "Error when setting gain for element: %s: %s",
					gain, err.Error())
				return
			}
			eltGain := sdr.GetGainElement(direction, channel, gain)
			log.Logf(logger.Info, "Gain for element %s is set to %.0f db\n", gain, eltGain)
			err = sdr.SetGainElement(direction, channel, gain, 0.0)
			log.Logf(logger.Info, "Have reset gain for element: %s to 0 db\n", gain)
		}
		log.Log(logger.Info, "Setting overall gain to 25 db\n")
		err = sdr.SetGain(direction, channel, 25.0)
		if err != nil {
			log.Logf(logger.Error, "Error when setting gain: %s", err.Error())
			return
		}
		var gainMsg strings.Builder
		gainMsg.WriteString(fmt.Sprintf("Overall gain set to: %.0f db\n", sdr.GetGain(direction, channel)))
		for _, gain := range gains {
			eltGain := sdr.GetGainElement(direction, channel, gain)
			gainMsg.WriteString(fmt.Sprintf("         %s gain is %.0f db\n", gain, eltGain))
		}
		log.Log(logger.Info, gainMsg.String())
	}
}

func exerciseSampleRate(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	sampleRanges := sdr.GetSampleRateRange(direction, channel)
	if len(sampleRanges) == 0 {
		log.Logf(logger.Info, "Channel#%d has no sample rate ranges\n", channel)
	} else {
		var sMsg strings.Builder
		sMsg.WriteString(fmt.Sprintf("Sample Rate Ranges for Channel#%d:\n", channel))
		for _, rng := range sampleRanges {
			sMsg.WriteString(fmt.Sprintf("         %v\n", rng))
		}
		log.Log(logger.Info, sMsg.String())
		log.Logf(logger.Info, "Channel#%d Sample Rate: %.0f\n", channel,
			sdr.GetSampleRate(direction, channel))
	}
	log.Log(logger.Info, "Setting sample rate to 1024000\n")
	err := sdr.SetSampleRate(direction, channel, 1024000.0)
	if err != nil {
		log.Logf(logger.Error, "Error while setting sample rate: %s", err.Error())
	}
	log.Logf(logger.Info, "Sample Rate is now %.0f\n",
		sdr.GetSampleRate(direction, channel))
}

func exerciseFrequencies(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	args := sdr.GetFrequencyArgsInfo(direction, channel)
	var aMsg strings.Builder
	if len(args) == 0 {
		log.Log(logger.Info, "Frequency Args Info: none\n")
	} else {
		aMsg.WriteString("Frequency Args Info:\n")
		for _, arg := range args {
			aMsg.WriteString(fmt.Sprintf("         %v\n", arg))
		}
		log.Log(logger.Info, aMsg.String())
	}
	freqRanges := sdr.GetFrequencyRange(direction, channel)
	if len(freqRanges) == 0 {
		log.Logf(logger.Info, "Channel#%d has no frequency ranges\n", channel)
		return
	}
	var fMsg strings.Builder
	fMsg.WriteString(fmt.Sprintf("Channel#%d Frequency Ranges:\n", channel))
	for _, fRange := range freqRanges {
		fMsg.WriteString(fmt.Sprintf("         %v\n", fRange))
	}
	log.Log(logger.Info, fMsg.String())

	tuneableElts := sdr.ListFrequencies(direction, channel)
	if len(tuneableElts) == 0 {
		log.Logf(logger.Info, "Channel#%d has no tuneable elements\n", channel)
	} else {
		var tMsg strings.Builder
		tMsg.WriteString(fmt.Sprintf("Tuneable elements for channel#%d:\n", channel))
		for _, elt := range tuneableElts {
			tMsg.WriteString(fmt.Sprintf("         %s\n", elt))
			comp := sdr.GetFrequencyComponent(direction, channel, elt)
			if elt == "CORR" {
				tMsg.WriteString(fmt.Sprintf("             Correction: %.0f PPM\n", comp))
				tMsg.WriteString("             Setting CORR to 50 PPM\n")
				err := sdr.SetFrequencyComponent(direction, channel, "CORR", 50.0, map[string]string{})
				if err != nil {
					log.Logf(logger.Error,
						"Error encountered setting CORR: %s", err.Error())
					return
				}
				tMsg.WriteString(fmt.Sprintf("             CORR now: %.0f PPM\n",
					sdr.GetFrequencyComponent(direction, channel, "CORR")))
			} else {
				tMsg.WriteString(fmt.Sprintf("             Frequency: %.0f Hz\n", comp))
				tMsg.WriteString("             Center Frequency set to 50 MHz\n")
				err := sdr.SetFrequencyComponent(direction, channel, elt, 50000000.0, map[string]string{})
				if err != nil {
					log.Logf(logger.Error,
						"Error encountered setting component %s frequency: %s\n",
						elt, err.Error())
					return
				}
				tMsg.WriteString(fmt.Sprintf("             Center Frequency now: %.0f Hz\n",
					sdr.GetFrequencyComponent(direction, channel, elt)))
			}
			rngs := sdr.GetFrequencyRangeComponent(direction, channel, elt)
			for _, rng := range rngs {
				tMsg.WriteString(fmt.Sprintf("             Range: %v\n", rng))
			}
		}
		log.Log(logger.Info, tMsg.String())

		log.Logf(logger.Info, "Current Center Frequency is: %.0f\n",
			sdr.GetFrequency(direction, channel))
		log.Log(logger.Info, "Setting Center Frequency to 75 MHz\n")
		err := sdr.SetFrequency(direction, channel, 75000000.0, map[string]string{})
		if err != nil {
			log.Logf(logger.Error, "Error encountered setting center frequency: %s",
				err.Error())
			return
		}
		log.Logf(logger.Info, "Center Frequency is now: %.0f\n",
			sdr.GetFrequency(direction, channel))
		var cMsg strings.Builder
		cMsg.WriteString("Component frequencies:\n")
		for _, elt := range tuneableElts {
			freq := sdr.GetFrequencyComponent(direction, channel, elt)
			if elt == "CORR" {
				cMsg.WriteString(fmt.Sprintf("         CORR: %.0f PPM\n", freq))
			} else {
				cMsg.WriteString(fmt.Sprintf("         %s: %.0f\n", elt, freq))
			}
		}
		log.Log(logger.Info, cMsg.String())
	}
}

func logStreamFormatsAndInfo(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	formats := sdr.GetStreamFormats(direction, channel)
	if len(formats) == 0 {
		log.Logf(logger.Info, "Channel#%d has no stream formats\n", channel)
		return
	}
	var fMsg strings.Builder
	fMsg.WriteString(fmt.Sprintf("Channel#%d stream formats:\n", channel))
	for _, format := range formats {
		fMsg.WriteString(fmt.Sprintf("         %s\n", format))
	}
	log.Log(logger.Info, fMsg.String())

	format, fullScale := sdr.GetNativeStreamFormat(direction, channel)
	log.Logf(logger.Info, "Native stream format: %s\n         fullScale: %f\n",
		format, fullScale)

	args := sdr.GetStreamArgsInfo(direction, channel)
	if len(args) == 0 {
		log.Logf(logger.Info, "Stream args info for channel#%d: none\nS", channel)
		return
	}
	var aMsg strings.Builder
	aMsg.WriteString(fmt.Sprintf("Stream Args Info for channel#%d:\n", channel))
	for _, arg := range args {
		aMsg.WriteString(fmt.Sprintf("         %v\n", arg))
	}
	log.Log(logger.Info, aMsg.String())
}

func exerciseFrontend(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	available := sdr.HasDCOffsetMode(direction, channel)
	if !available {
		log.Logf(logger.Info, "Channel#%d does not support stream auto DC Correction available\n",
			channel)
	} else {
		offsetMode := sdr.GetDCOffsetMode(direction, channel)
		log.Logf(logger.Info, "Channel#%d supports auto DC Offset correction: %v\n",
			channel, offsetMode)
	}
	available = sdr.HasDCOffset(direction, channel)
	if !available {
		log.Logf(logger.Info, "Channel#%d does not support frontend DC offset correction\n",
			channel)
	} else {
		log.Logf(logger.Info, "Channel#%d supports frontend DC offset correction\n",
			channel)
		offsetI, offsetQ, err := sdr.GetDCOffset(direction, channel)
		if err != nil {
			log.Logf(logger.Error, "Error encountered retrieving stream DCOffset: %s\n",
				err.Error())
			return
		}
		log.Logf(logger.Info, "Channel#%d stream DC offset relative correction: I: %f, Q: %f\n",
			channel, offsetI, offsetQ)
		log.Log(logger.Info, "Setting stream DC offset to 0.1, 0.1\n")
		err = sdr.SetDCOffset(direction, channel, 0.1, 0.1)
		if err != nil {
			log.Logf(logger.Error, "Error encountered setting stream DC offset: %s\n", err.Error())
			return
		}
		offsetI, offsetQ, err = sdr.GetDCOffset(direction, channel)
		log.Logf(logger.Info, "Channel#%d stream DC offset now: I: %f, Q: %f\n",
			channel, offsetI, offsetQ)
	}
	available = sdr.HasIQBalance(direction, channel)
	if !available {
		log.Logf(logger.Info, "Channel#%d does not support IQ Balance\n", channel)
	} else {
		log.Logf(logger.Info, "Channel#%d supports IQ Balance\n", channel)
		I, Q, err := sdr.GetIQBalance(direction, channel)
		if err != nil {
			log.Logf(logger.Info, "Error encountered getting I, Q balance values: %s\n",
				err.Error())
			return
		}
		log.Logf(logger.Info, "Channel#%d I/Q Balance: I: %f, Q: %f\n",
			channel, I, Q)
		log.Log(logger.Info, "Setting I/Q balance to 0.1, 0.1\n")
		err = sdr.SetIQBalance(direction, channel, 0.1, 0.1)
		if err != nil {
			log.Logf(logger.Info, "Error encountered setting I/Q balance: %s\n", err.Error())
			return
		}
		I, Q, err = sdr.GetIQBalance(direction, channel)
		if err != nil {
			log.Logf(logger.Info, "Error encountered getting I, Q balance values: %s\n",
				err.Error())
			return
		}
		log.Logf(logger.Info, "Channel#%d I/Q balance now set to I: %f, Q: %f\n",
			channel, I, Q)
	}
	available = sdr.HasFrequencyCorrection(direction, channel)
	if !available {
		log.Logf(logger.Info, "Channel#%d does not support frontend frequency correction\n",
			channel)
	} else {
		log.Logf(logger.Info, "Channel#%d supports frontend frequency correction\n",
			channel)
		correction := sdr.GetFrequencyCorrection(direction, channel)
		log.Logf(logger.Info, "Channel#%d frontend frequency correction is: %f\n",
			channel, correction)
		log.Log(logger.Info, "Setting frequency correction to 127 PPM\n")
		err := sdr.SetFrequencyCorrection(direction, channel, 127.0)
		if err != nil {
			log.Logf(logger.Error, "Error encountered setting frontend frequency correction: %s\n",
				err.Error())
			return
		}
		correction = sdr.GetFrequencyCorrection(direction, channel)
		log.Logf(logger.Info, "Frontend frequency correction is now %f PPM\n", correction)
	}
}

func exerciseChannelSensors(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	sensors := sdr.ListChannelSensors(direction, channel)
	if len(sensors) == 0 {
		log.Logf(logger.Info, "Channel#%d does not have any sensors\n", channel)
		return
	}
	var sMsg strings.Builder
	sMsg.WriteString(fmt.Sprintf("Channel#%d Sensors:\n", channel))
	for _, sensor := range sensors {
		sMsg.WriteString(fmt.Sprintf("         %s\n", sensor))
		args := sdr.GetChannelSensorInfo(direction, channel, sensor)
		sMsg.WriteString(fmt.Sprintf("            %v\n", args))
		sensorValue := sdr.ReadChannelSensor(direction, channel, sensor)
		sMsg.WriteString(fmt.Sprintf("            Current value: %s\n", sensorValue))
	}
	log.Log(logger.Info, sMsg.String())
}

func exerciseStream0(sdr *device.SDRDevice, direction device.Direction, channel uint, log *logger.Logger) {
	// set up stream params
	log.Log(logger.Debug, "Setting up stream params\n")
	log.Log(logger.Debug, "Setting sample rate to 1024e3\n")
	sdrErr := sdr.SetSampleRate(direction, channel, 1024e3)
	if sdrErr != nil {
		log.Logf(logger.Error, "Error encountered setting sample rate: %s\n", sdrErr.Error())
		return
	}
	log.Log(logger.Debug, "Setting center frequency to 99.9 MHz\n")
	sdrErr = sdr.SetFrequency(direction, channel, 99900000.0, map[string]string{})
	if sdrErr != nil {
		log.Logf(logger.Error, "Error encountered setting center frequency to 99.9 MHz: %s\n",
			sdrErr.Error())
		return
	}
	log.Log(logger.Debug, "Setting up the stream\n")
	stream, err := sdr.SetupSDRStreamCS8(direction, []uint{0}, nil)
	if err != nil {
		log.Logf(logger.Error, "Error encountered in SetupSDRStreamCS8 call: %s\n", err.Error())
		return
	}
	defer func() {
		log.Log(logger.Debug, "Closing the stream\n")
		err = stream.Close()
		if err != nil {
			log.Logf(logger.Error, "Error encountered closing the stream: %s\n", err.Error())
		}
		log.Log(logger.Debug, "Stream closed\n")

	}()

	log.Log(logger.Debug, "Activating the stream\n")
	err = stream.Activate(0, 0, 0)
	if err != nil {
		log.Logf(logger.Error, "Error encountered activating the stream: %s\n", err.Error())
		return
	}
	log.Log(logger.Debug, "Stream activated\n")

	defer func() {
		log.Log(logger.Debug, "Deactivating the stream\n")
		log.Log(logger.Debug, "Need to wait short while to ensure that stream.Deactivate does not hang\n")
		// Any amount of time seems sufficient in my testing. Chose this as a compromise.
		// Note: Deactivate never hangs when debugging.
		time.Sleep(10 * time.Millisecond)
		err = stream.Deactivate(0, 0)
		if err != nil {
			log.Logf(logger.Error, "Error encountered deactivating the stream: %s\n",
				err.Error())
			return
		}
		log.Log(logger.Debug, "Stream deactivated\n")
	}()

	mtu := stream.GetMTU()
	log.Logf(logger.Info, "Stream MTU: %d\n", mtu)

	numBuffers := stream.GetNumDirectAccessBuffers()
	log.Logf(logger.Info, "Number of Direct Access Buffers: %d\n", numBuffers)

	buffers := make([][]int8, 1)
	buffers[0] = make([]int8, 2*mtu)
	flags := make([]int, 1)

	for i := 0; i < 10; i++ {
		var rMsg strings.Builder
		rMsg.WriteString(fmt.Sprintf("Reading stream (%d)\n", i))
		timeNS, numElementsRead, err := stream.Read(buffers, uint(mtu), flags, 5000000)
		rMsg.WriteString(fmt.Sprintf("         timestamp: %d\n", timeNS))
		rMsg.WriteString(fmt.Sprintf("         numer of Elts read: %d\n", numElementsRead))
		rMsg.WriteString(fmt.Sprintf("         err: %v\n", err))
		rMsg.WriteString("         ")
		for j := 0; j < 7; j = j + 2 {
			rMsg.WriteString(fmt.Sprintf("[%d, %d]", buffers[0][2*j], buffers[0][2*j+1]))
		}
		rMsg.WriteString("\n         ...\n         ")
		for j := mtu - 1; j > mtu-9; j = j - 2 {
			rMsg.WriteString(fmt.Sprintf("[%d, %d]", buffers[0][2*j-1], buffers[0][2*j]))
		}
		rMsg.WriteString("\n")
		log.Log(logger.Info, rMsg.String())
	}
}
