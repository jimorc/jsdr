package jsdrgui

import (
	"fmt"
	"internal/settings"
	"internal/soapylogging"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

type modalPopUp struct {
	popUp *widget.PopUp
}

type loggingFileNameEntry struct {
	entry *widget.Entry
}

var loggingLevelSelect *widget.Select
var loggingPopUp = modalPopUp{}

// newSDRLoggerSettingsPopUp creates the logging modal popup.
// The return value is a pointer to the modal popup. This popup is displayed over the window specified in the
// calling parameter when popup.Show() is called.
// The popup is used to review and change logging parameters such as the name of the logging file and the the logging level.
func newSDRLoggerSettingsPopUp(win *fyne.Window) *widget.PopUp {
	loggingLevelLabel := widget.NewLabel("SDR Logging Level:")
	loggingLevelSelect = widget.NewSelect([]string{"Fatal", "Critical", "Error", "Warning", "Notice", "Info", "Debug",
		"Trace", "SSI"}, nil)
	loggingLevelSelect.SetSelectedIndex(int(atomic.LoadInt64(&settings.JsdrSettings.LoggingLevel) - 1))
	container := container.NewGridWithColumns(2, loggingLevelLabel, loggingLevelSelect,
		widget.NewButton("Reset", resetLoggingValues), widget.NewButton("Accept", acceptChanges))
	loggingPopUp.popUp = widget.NewModalPopUp(container, (*win).Canvas())
	return loggingPopUp.popUp
}

// acceptChanges processes clicks on the "Accept" button.
//
// It saves the logging level to the JsdrSettings object, and, if the log file name is valid and not the same as the
// value in JsdrSettings, then it saves the filename, renames the logging file, and closes the popup window.
func acceptChanges() {
	level := int64(loggingLevelSelect.SelectedIndex() + 1)
	atomic.StoreInt64(&settings.JsdrSettings.LoggingLevel, level)
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("acceptChanges - set logging level to %v",
		soapylogging.LoggingLevelAsString(sdrlogger.SDRLogLevel(level))))
	loggingPopUp.closeLoggingPopUp()
}

// resetLoggingValues resets the loggingFileName entry and the loggingLevelSelect to the values in JsdrSettings.
func resetLoggingValues() {
	loggingLevelSelect.SetSelectedIndex(int(atomic.LoadInt64(&settings.JsdrSettings.LoggingLevel) - 1))
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Logging level reset to: %v",
		loggingLevelSelect.Selected))
}

// closeLoggingPopUp closes the logging popup window.
func (modPopUp *modalPopUp) closeLoggingPopUp() {
	modPopUp.popUp.Hide()
}
