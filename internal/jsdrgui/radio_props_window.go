package jsdrgui

import (
	"fmt"
	"internal/settings"
	"internal/soapydevice"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

type ppmCorrection struct {
	entry *widget.Entry
}

var samplingModeSettings []string

var radioSelect = widget.NewSelect([]string{""}, radioWindow.radioSelected)
var sampleRates = widget.NewSelect([]string{""}, radioWindow.sampleRateSelected)
var antennaSelect = widget.NewSelect([]string{""}, radioWindow.antennaSelected)
var samplingModeSelect = widget.NewSelect([]string{""}, radioWindow.samplingModeSelected)

var frequencyCorrection = &ppmCorrection{}
var ppm = binding.NewInt()
var ppmCorr = binding.IntToString(ppm)

var layoutWidth float32 = 450.0

// radioWin is the window containing the radio settings.
var radioWindow *actionWindow = nil

// newRadioWindow creates the radio popup window
// The return value is a pointer to the radioWindow struct. The window is displayed over the window specified in the
// calling parameter when window.Show() is called.
// The window is used to select an SDR device and some of its parameters.
// If there are no SDRs attached to the computer, an information message is displayed, and nil is returned
func newRadioWindow(parent *fyne.Window) *actionWindow {
	radioWindow = &actionWindow{}
	radioWindow.window = SdrApp.NewWindow("Radio Properties")
	radioLabel := widget.NewLabel("Radio:")
	radioLabel.Alignment = fyne.TextAlignTrailing
	sampleRateLabel := widget.NewLabel("A/D Sample Rate:")
	sampleRateLabel.Alignment = fyne.TextAlignTrailing
	antennaLabel := widget.NewLabel("Antenna:")
	antennaLabel.Alignment = fyne.TextAlignTrailing
	samplingModeLabel := widget.NewLabel("Sampling Mode:")
	samplingModeLabel.Alignment = fyne.TextAlignTrailing
	frequencyCorrectionLabel := widget.NewLabel("PPM Correction:")
	frequencyCorrectionLabel.Alignment = fyne.TextAlignTrailing
	frequencyCorrection.entry = widget.NewEntryWithData(ppmCorr)

	formContainer := &fyne.Container{
		Objects: []fyne.CanvasObject{radioLabel, radioSelect, sampleRateLabel, sampleRates, antennaLabel, antennaSelect,
			samplingModeLabel, samplingModeSelect, frequencyCorrectionLabel, frequencyCorrection.entry},
	}
	layout := layout.NewFormLayout()
	layout.Layout(formContainer.Objects, fyne.NewSize(layoutWidth, 150))

	accept := widget.NewButton("Accept", radioAcceptChanges)
	cancel := widget.NewButton("Cancel", radioCancelChanges)

	buttonBar := container.NewHBox()
	buttonBar.Add(cancel)
	buttonBar.Add(accept)
	buttonBox := container.NewBorder(nil, nil, nil, buttonBar)

	cont := container.NewBorder(formContainer, buttonBox, nil, nil)
	radioWindow.window.SetContent(cont)
	radioWindow.window.SetOnClosed(closeRadioWindow)

	radioWindow.window.Resize(fyne.NewSize(layoutWidth+3*theme.Padding(), 250))

	radios := device.Enumerate(nil)
	if len(radios) == 0 {
		sdrlogger.Logf(sdrlogger.Trace, "No radios found")
		dialog.ShowInformation("No Radios", "No SDR radios found.\nPlease attach an SDR and click\nthe Radio toolbar item again.",
			*parent)
		return nil
	}
	var labels []string
	for _, radio := range radios {
		sdrlogger.Logf(sdrlogger.Trace, "Found SDR: %v", radio["label"])
		labels = append(labels, radio["label"])
	}
	radioSelect.SetOptions(labels)
	if len(radios) == 1 {
		radioSelect.SetSelectedIndex(0)
	} else if len(settings.JsdrSettings.Sdr) > 0 {
		radioSelect.SetSelected(settings.JsdrSettings.Sdr)
	}
	return radioWindow
}

// acceptChanges processes clicks on the "Accept" button.
func radioAcceptChanges() {
	sdrlogger.Log(sdrlogger.Trace, "In radioAcceptChanges")
	if radioSelect.Selected != settings.JsdrSettings.Sdr {
		sdrlogger.Logf(sdrlogger.Trace, fmt.Sprintf("JsdrSettings.Sdr set to %v", radioSelect.Selected))
		settings.JsdrSettings.Sdr = radioSelect.Selected
	}
	if frequencyCorrection.entry.Validate() != nil {
		corr := frequencyCorrection.entry.Text
		sdrlogger.Log(sdrlogger.Error, fmt.Sprintf("Invalid frequencyCorrection value: '%v' when Accept button pressed"+
			" in Hardware dialog", corr))
		dialog.ShowInformation("Invalid PPM", "Invalid PPM Correction.\nMust be an integer value.\n"+
			"Either correct the value and click 'Accept' again,\nor click 'Cancel",
			actionWin.window)
		return
	}
	actionWin.window.Close()
}

func radioCancelChanges() {
	actionWin.window.Close()
}

// resetRadioValues resets the radio entry.
func rescanRadioValues() {
	radioSelect.SetSelectedIndex(0)
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Radio set to: %v",
		radioSelect.Selected))
}

// closeRadioWindow closes the radio window.
func closeRadioWindow() {
	actionWin = nil

}

// radioSelected retrieves SDR properties for display when an SDR is selected.
func (radioWin *actionWindow) radioSelected(sdr string) {
	sdrlogger.Logf(sdrlogger.Trace, "SDR: %v selected", sdr)
	deviceArgs := make(map[string]string, 1)
	deviceArgs = map[string]string{
		"label": sdr,
	}

	if soapydevice.Radio != nil {
		soapydevice.Radio.Unmake()
	}

	dev, err := soapydevice.Make(deviceArgs)
	if err != nil {
		dialog.ShowInformation("SDR Not Found", fmt.Sprintf("An error has occurred.\nCannnot access the selected SDR:\n%v", err),
			radioWindow.window)
	}
	soapydevice.Radio = dev
	soapydevice.Radio.GetSampleRateRange()
	sampleRates.SetOptions(soapydevice.Radio.SampleRates)

	soapydevice.Radio.GetListOfAntennas()
	antennaSelect.SetOptions(soapydevice.Radio.Antennas)
	if len(soapydevice.Radio.Antennas) == 1 {
		antennaSelect.SetSelectedIndex(0)
	}

	setValue, samplingNames := soapydevice.Radio.GetSamplingModeNames()
	samplingModeSelect.SetOptions(samplingNames)
	index, err := strconv.ParseInt(setValue, 10, 32)
	samplingModeSelect.SetSelectedIndex(int(index))

	frequencyComponents := soapydevice.Radio.ListFrequencies()
	for _, component := range frequencyComponents {
		if component == "CORR" {
			ppm.Set(int(soapydevice.Radio.GetFrequencyCorrection()))
		}
	}
}

func (radioWin *actionWindow) sampleRateSelected(rate string) {
	if err := soapydevice.Radio.SetSampleRate(rate); err != nil {
		dialog.ShowInformation("Sample Rate Error", fmt.Sprintf("Error encountered attempting to set sample rate %v:", rate),
			radioWin.window)
		sdrlogger.Logf(sdrlogger.Trace, "Sample rate of %v selected", rate)
	}
}

func (radioWin *actionWindow) antennaSelected(antenna string) {
	if err := soapydevice.Radio.SetAntenna(antenna); err != nil {
		dialog.ShowInformation("Set Antenna Error",
			fmt.Sprintf("Could not set antenna: %v\n%v\nCheck the SDR and possibly select a different antenna.", antenna, err),
			radioWin.window)
		return
	}
	sdrlogger.Logf(sdrlogger.Trace, "Antenna %v selected", antenna)
}

func (radioWin *actionWindow) samplingModeSelected(mode string) {
	sdrlogger.Logf(sdrlogger.Trace, "Sampling mode %v selected", mode)
}
