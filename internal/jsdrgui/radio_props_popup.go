package jsdrgui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

type radioPopUp struct {
	popUp *widget.PopUp
}

type radioEntry struct {
	entry *widget.Entry
}

var radioSelect *widget.Select
var radioPopup = radioPopUp{}

// newSDRLoggerSettingsPopUp creates the logging modal popup.
// The return value is a pointer to the modal popup. This popup is displayed over the window specified in the
// calling parameter when popup.Show() is called.
// The popup is used to review and change logging parameters such as the name of the logging file and the the logging level.
func newRadioPopUp(win *fyne.Window) *widget.PopUp {
	radioLabel := widget.NewLabel("Radio:")
	radioSelect = widget.NewSelect([]string{"Radio 0", "Radio 1", "Radio 2"}, nil)
	radioSelect.SetSelectedIndex(0)
	container := container.NewGridWithColumns(2, radioLabel, radioSelect,
		widget.NewButton("Reset", resetRadioValues), widget.NewButton("Accept", radioAcceptChanges))
	radioPopup.popUp = widget.NewModalPopUp(container, (*win).Canvas())
	return radioPopup.popUp
}

// acceptChanges processes clicks on the "Accept" button.
//
// It saves the logging level to the JsdrSettings object, and, if the log file name is valid and not the same as the
// value in JsdrSettings, then it saves the filename, renames the logging file, and closes the popup window.
func radioAcceptChanges() {
	sdrlogger.Log(sdrlogger.Trace, "In radioAcceptChanges")
	radioPopup.closeRadioPopUp()
}

// resetRadioValues resets the radio entry.
func resetRadioValues() {
	radioSelect.SetSelectedIndex(0)
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Radio set to: %v",
		radioSelect.Selected))
}

// closeLoggingPopUp closes the logging popup window.
func (rPopUp *radioPopUp) closeRadioPopUp() {
	sdrlogger.Log(sdrlogger.Trace, "Closing radio popup")
	rPopUp.popUp.Hide()
}
