package jsdrgui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/pothosware/go-soapy-sdr/pkg/device"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

type radioPopUp struct {
	popUp *widget.PopUp
	// parent *fyne.Window
}

type radioEntry struct {
	entry *widget.Entry
}

var radioSelect *widget.Select
var radioPopup = radioPopUp{}

// newRadioPopUp creates the logging modal popup.
// The return value is a pointer to the modal popup. This popup is displayed over the window specified in the
// calling parameter when popup.Show() is called.
// The popup is used to select an SDR device and some of its parame.ters.
// If there are no SDRs attached to the computer, an information message is displayed, and nil is returned
func newRadioPopUp(win *fyne.Window) *widget.PopUp {
	radioLabel := widget.NewLabel("Radio:")
	radioSelect = widget.NewSelect([]string{""}, nil)
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
	radioSelect.SetSelectedIndex(0)
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
