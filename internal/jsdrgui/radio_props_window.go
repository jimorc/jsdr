package jsdrgui

import (
	"fmt"
	"internal/settings"
	"internal/soapydevice"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

type radioEntry struct {
	entry *widget.Entry
}

var radioSelect = widget.NewSelect([]string{""}, radioWindow.radioSelected)

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
	radioLabel := widget.NewLabel("Radio")
	container := container.NewGridWithColumns(2, radioLabel, radioSelect,
		widget.NewButton("Rescan", rescanRadioValues), widget.NewButton("Accept", radioAcceptChanges))
	radioWindow.window.SetContent(container)
	radioWindow.window.SetOnClosed(closeRadioWindow)

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
}
