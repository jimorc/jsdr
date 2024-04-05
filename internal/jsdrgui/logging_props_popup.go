package jsdrgui

import (
	"fmt"
	"internal/settings"
	"internal/soapylogging"
	"os"

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

var loggingFileName *loggingFileNameEntry
var loggingLevelSelect *widget.Select

func newLoggingFileNameEntry() *loggingFileNameEntry {
	fileNameEntry := loggingFileNameEntry{entry: widget.NewEntry()}
	fileNameEntry.entry.SetText(settings.JsdrSettings.Logging.LoggingFile)
	fileNameEntry.entry.OnSubmitted = loggingFileNameSubmitted
	fileNameEntry.entry.OnChanged = fileNameEntry.loggingFileNameChanged
	fileNameEntry.entry.Validator = fileNameEntry.validateLoggingFileName
	return &fileNameEntry
}

// newSDRLoggerSettingsPopUp creates the logging modal popup.
// The return value is a pointer to the modal popup. This popup is displayed over the window specified in the
// calling parameter when popup.Show() is called.
// The popup is used to review and change logging parameters such as the name of the logging file and the the logging level.
func newSDRLoggerSettingsPopUp(win *fyne.Window) *widget.PopUp {
	loggingFileName = newLoggingFileNameEntry()

	loggingFileLabel := widget.NewLabel("SDR Logging File Name:")
	loggingLevelLabel := widget.NewLabel("SDR Logging Level:")
	loggingLevelSelect = widget.NewSelect([]string{"Fatal", "Critical", "Error", "Warning", "Notice", "Info", "Debug",
		"Trace", "SSI"}, nil)
	loggingLevelSelect.SetSelectedIndex(int(settings.JsdrSettings.Logging.LoggingLevel - 1))
	var loggingPopUp modalPopUp
	container := container.NewGridWithColumns(2, loggingFileLabel, loggingFileName.entry, loggingLevelLabel, loggingLevelSelect,
		widget.NewButton("Reset", resetLoggingValues), widget.NewButton("Close", loggingPopUp.closeLoggingPopUp))
	loggingPopUp.popUp = widget.NewModalPopUp(container, (*win).Canvas())
	return loggingPopUp.popUp
}

// resetLoggingValues resets the loggingFileName entry and the loggingLevelSelect to the values in JsdrSettings.
func resetLoggingValues() {
	loggingFileName.entry.SetText(settings.JsdrSettings.Logging.LoggingFile)
	loggingLevelSelect.SetSelectedIndex(int(settings.JsdrSettings.Logging.LoggingLevel - 1))
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Logging values reset to: %v, %v",
		loggingFileName.entry.Text, loggingLevelSelect.Selected))
}

func (entry *loggingFileNameEntry) loggingFileNameChanged(fileName string) {
	fmt.Printf("In loggingFileNameChanged. FileName is: %v\n", fileName)
}

func loggingFileNameSubmitted(filename string) {
	fmt.Printf("Submitted: File name: %v\n", filename)
}

// validateLoggingFileName validates the filename in loggingFileNameEntry
//
// Filename is valid if:
// 1. It is the same as the file name in settings.JsdrSettings.Logging.LoggingFile; or,
// 2. It exists and can be opened for writing; or,
// 3. It can be created, and therefore opened for writing.
func (entry *loggingFileNameEntry) validateLoggingFileName(filename string) error {
	fileName := entry.entry.Text
	if fileName == settings.JsdrSettings.Logging.LoggingFile {
		sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("validateLoggingFileName - File: %v matches settings", filename))
		return nil
	}
	soapylogging.SoapyLoggingMutex.Lock()
	defer soapylogging.SoapyLoggingMutex.Unlock()
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("validateLoggingFileName - File: %v cannot be opened: %v", filename, err))
		file, err = os.Create(fileName)
		if err != nil {
			sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("validateLoggingFileName - File: %v cannot be created: %v", filename, err))
			return err
		}
		file.Close()
		os.Remove(fileName)
		sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("validateLoggingFileName - File: %v can be created", filename))
		return nil
	}
	file.Close()
	sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("validateLoggingFileName - File: %v exists and can be opened for writing", filename))
	return nil
}

// closeLoggingPopUp closes the logging popup window.
func (modPopUp *modalPopUp) closeLoggingPopUp() {
	modPopUp.popUp.Hide()
}
