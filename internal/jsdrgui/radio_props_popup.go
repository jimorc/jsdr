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

type radioPopUp struct {
	popUp  *widget.PopUp
	parent *fyne.Window
}

type radioEntry struct {
	entry *widget.Entry
}

var radioSelect *widget.Select
var radioPopup = radioPopUp{}

// newRadioPopUp creates the logging modal popup.
// The return value is a pointer to the modal popup. This popup is displayed over the window specified in the
// calling parameter when popup.Show() is called.
// The popup is used to select an SDR device and some of its parameters.
// If there are no SDRs attached to the computer, an information message is displayed, and nil is returned
func newRadioPopUp(win *fyne.Window) *widget.PopUp {
	radioPopup.parent = win
	radioLabel := widget.NewLabel("Radio:")
	radioSelect = widget.NewSelect([]string{""}, radioPopup.radioSelected)
	container := container.NewGridWithColumns(2, radioLabel, radioSelect,
		widget.NewButton("Rescan", rescanRadioValues), widget.NewButton("Accept", radioAcceptChanges))
	radioPopup.popUp = widget.NewModalPopUp(container, (*win).Canvas())

	radios := device.Enumerate(nil)
	if len(radios) == 0 {
		sdrlogger.Logf(sdrlogger.Trace, "No radios found")
		dialog.ShowInformation("No Radios", "No SDR radios found.\nPlease attach an SDR and click\nthe Radio toolbar item again.",
			*win)
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
	return radioPopup.popUp
}

// acceptChanges processes clicks on the "Accept" button.
func radioAcceptChanges() {
	sdrlogger.Log(sdrlogger.Trace, "In radioAcceptChanges")
	radioPopup.closeRadioPopUp()
}

// resetRadioValues resets the radio entry.
func rescanRadioValues() {
	radioSelect.SetSelectedIndex(0)
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Radio set to: %v",
		radioSelect.Selected))
}

// closeLoggingPopUp closes the logging popup window.
func (rPopUp *radioPopUp) closeRadioPopUp() {
	sdrlogger.Log(sdrlogger.Trace, "Closing radio popup")
	rPopUp.popUp.Hide()
}

// radioSelected retrieves SDR properties for display when an SDR is selected.
func (rPopUp *radioPopUp) radioSelected(sdr string) {
	sdrlogger.Logf(sdrlogger.Trace, "SDR: %v selected", sdr)
	if sdr != settings.JsdrSettings.Sdr {
		settings.JsdrSettings.Sdr = sdr
	}
	deviceArgs := make([]map[string]string, 1)
	deviceArgs[0] = map[string]string{
		"label": sdr,
	}
	_, err := device.MakeList(deviceArgs)
	if err != nil {
		sdrlogger.Logf(sdrlogger.Error, "Error retrieving the selected SDR: %v", err)
		rPopUp.popUp.Hide()
		dialog.ShowInformation("SDR Not Found", fmt.Sprintf("An error has occurred.\nCannnot access the selected SDR:\n%v", err),
			*rPopUp.parent)

	}

}
