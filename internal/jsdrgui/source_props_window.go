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

var sourceSelect = widget.NewSelect([]string{""}, sourceSelected)
var sampleRates = widget.NewSelect([]string{""}, sampleRateSelected)
var antennaSelect = widget.NewSelect([]string{""}, antennaSelected)
var samplingModeSelect = widget.NewSelect([]string{""}, samplingModeSelected)

var frequencyCorrection = &ppmCorrection{}
var ppm = binding.NewInt()
var ppmCorr = binding.IntToString(ppm)

var layoutWidth float32 = 450.0

// newSourceWindow creates the source popup window
// The return value is a pointer to the sourceWindow struct. The window is displayed over the window specified in the
// calling parameter when window.Show() is called.
// The window is used to select an SDR device or other source and some of its parameters.
// If there are no SDRs attached to the computer, an information message is displayed, and nil is returned
func newSourceWindow(parent *fyne.Window) *actionWindow {
	sourceWindow := &actionWindow{}
	sourceWindow.window = SdrApp.NewWindow("Source Properties")
	sourceLabel := widget.NewLabel("Source:")
	sourceLabel.Alignment = fyne.TextAlignTrailing
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
		Objects: []fyne.CanvasObject{sourceLabel, sourceSelect, sampleRateLabel, sampleRates, antennaLabel, antennaSelect,
			samplingModeLabel, samplingModeSelect, frequencyCorrectionLabel, frequencyCorrection.entry},
	}
	layout := layout.NewFormLayout()
	layout.Layout(formContainer.Objects, fyne.NewSize(layoutWidth, 150))

	accept := widget.NewButton("Accept", sourceAcceptChanges)
	cancel := widget.NewButton("Cancel", sourceCancelChanges)

	buttonBar := container.NewHBox()
	buttonBar.Add(cancel)
	buttonBar.Add(accept)
	buttonBox := container.NewBorder(nil, nil, nil, buttonBar)

	cont := container.NewBorder(formContainer, buttonBox, nil, nil)
	sourceWindow.window.SetContent(cont)
	sourceWindow.window.SetOnClosed(closeSourceWindow)

	sourceWindow.window.Resize(fyne.NewSize(layoutWidth+3*theme.Padding(), 250))

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
	sourceSelect.SetOptions(labels)
	if len(radios) == 1 {
		sourceSelect.SetSelectedIndex(0)
	} else if len(settings.JsdrSettings.Sdr) > 0 {
		sourceSelect.SetSelected(settings.JsdrSettings.Sdr)
	}
	return sourceWindow
}

// sourceAcceptChanges processes clicks on the "Accept" button.
func sourceAcceptChanges() {
	sdrlogger.Log(sdrlogger.Trace, "In radioAcceptChanges")
	if sourceSelect.Selected != settings.JsdrSettings.Sdr {
		sdrlogger.Logf(sdrlogger.Trace, fmt.Sprintf("JsdrSettings.Sdr set to %v", sourceSelect.Selected))
		settings.JsdrSettings.Sdr = sourceSelect.Selected
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

func sourceCancelChanges() {
	actionWin.window.Close()
}

// resetSourceValues resets the source entry to the first source.
func resetSourceValues() {
	sourceSelect.SetSelectedIndex(0)
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Radio set to: %v",
		sourceSelect.Selected))
}

// closeSourceWindow closes the source window.
func closeSourceWindow() {
	actionWin = nil

}

// sourceSelected retrieves SDR properties for display when an SDR is selected.
func sourceSelected(sdr string) {
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
			actionWin.window)
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

func sampleRateSelected(rate string) {
	if err := soapydevice.Radio.SetSampleRate(rate); err != nil {
		dialog.ShowInformation("Sample Rate Error", fmt.Sprintf("Error encountered attempting to set sample rate %v:", rate),
			actionWin.window)
		sdrlogger.Logf(sdrlogger.Trace, "Sample rate of %v selected", rate)
	}
}

func antennaSelected(antenna string) {
	if err := soapydevice.Radio.SetAntenna(antenna); err != nil {
		dialog.ShowInformation("Set Antenna Error",
			fmt.Sprintf("Could not set antenna: %v\n%v\nCheck the SDR and possibly select a different antenna.", antenna, err),
			actionWin.window)
		return
	}
	sdrlogger.Logf(sdrlogger.Trace, "Antenna %v selected", antenna)
}

func samplingModeSelected(mode string) {
	sdrlogger.Logf(sdrlogger.Trace, "Sampling mode %v selected", mode)
}
