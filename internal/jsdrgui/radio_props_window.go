package jsdrgui

import (
	"fmt"
	"internal/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

type radioWindow struct {
	Window fyne.Window
}

type radioEntry struct {
	entry *widget.Entry
}

var radioSelect = widget.NewSelect([]string{""}, radioWin.radioSelected)

// radioWin is the window containing the radio settings.
var radioWin *radioWindow = nil

// newRadioWindow creates the radio popup window
// The return value is a pointer to the radioWindow struct. The window is displayed over the window specified in the
// calling parameter when window.Show() is called.
// The window is used to select an SDR device and some of its parameters.
// If there are no SDRs attached to the computer, an information message is displayed, and nil is returned
func newRadioWindow(parent *fyne.Window) *radioWindow {
	radioWin = &radioWindow{}
	radioWin.Window = SdrApp.NewWindow("Radio Properties")
	radioLabel := widget.NewLabel("Radio")
	container := container.NewGridWithColumns(2, radioLabel, radioSelect,
		widget.NewButton("Rescan", rescanRadioValues), widget.NewButton("Accept", radioAcceptChanges))
	radioWin.Window.SetContent(container)
	radioWin.Window.SetOnClosed(closeRadioWindow)

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
	return radioWin
}

// acceptChanges processes clicks on the "Accept" button.
func radioAcceptChanges() {
	sdrlogger.Log(sdrlogger.Trace, "In radioAcceptChanges")
	if radioSelect.Selected != settings.JsdrSettings.Sdr {
		sdrlogger.Logf(sdrlogger.Trace, fmt.Sprintf("JsdrSettings.Sdr set to %v", radioSelect.Selected))
		settings.JsdrSettings.Sdr = radioSelect.Selected
	}
	radioWin.Window.Close()
}

// resetRadioValues resets the radio entry.
func rescanRadioValues() {
	radioSelect.SetSelectedIndex(0)
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Radio set to: %v",
		radioSelect.Selected))
}

// closeRadioWindow closes the radio window.
func closeRadioWindow() {
	radioWin = nil
	enableMainToolbar()
}

// radioSelected retrieves SDR properties for display when an SDR is selected.
func (radioWin *radioWindow) radioSelected(sdr string) {
	sdrlogger.Logf(sdrlogger.Trace, "SDR: %v selected", sdr)
	deviceArgs := make([]map[string]string, 1)
	deviceArgs[0] = map[string]string{
		"label": sdr,
	}
	devs, err := device.MakeList(deviceArgs)
	if err != nil {
		sdrlogger.Logf(sdrlogger.Error, "Error retrieving the selected SDR: %v", err)
		radioWin.Window.Hide()
		dialog.ShowInformation("SDR Not Found", fmt.Sprintf("An error has occurred.\nCannnot access the selected SDR:\n%v", err),
			radioWin.Window)
	}
	if len(devs) > 1 {
		sdrlogger.Logf(sdrlogger.Error, fmt.Sprintf("More than one SDR retrieved for the selected SDR: %v", deviceArgs[0]["label"]))
		radioWin.Window.Hide()
		dialog.ShowInformation("Multiple SDRs Found",
			"More than one SDR retrieved for the selected item.\nSee documentation for information about how SDRs are "+
				"distinguished.\nIf this does not explain the problem,\nfile a bug report and include the contents of the jsdr.log file.",
			radioWin.Window)
	}

}
