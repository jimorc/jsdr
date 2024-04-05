package jsdrgui

import (
	"fmt"
	"internal/settings"
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
	loggingFileName := newLoggingFileNameEntry()

	loggingFileLabel := widget.NewLabel("SDR Logging File Name:")
	loggingLevelLabel := widget.NewLabel("SDR Logging Level:")
	loggingLevelSelect := widget.NewSelect([]string{"Fatal", "Critical", "Error", "Warning", "Notice", "Info", "Debug",
		"Trace", "SSI"}, loggingLevelSelectChanged)
	var loggingPopUp modalPopUp
	container := container.NewGridWithColumns(2, loggingFileLabel, loggingFileName.entry, loggingLevelLabel, loggingLevelSelect,
		widget.NewLabel(""), widget.NewButton("Close", loggingPopUp.closeLoggingPopUp))
	loggingPopUp.popUp = widget.NewModalPopUp(container, (*win).Canvas())
	return loggingPopUp.popUp
}

func (entry *loggingFileNameEntry) loggingFileNameChanged(fileName string) {
	fmt.Printf("In loggingFileNameChanged. FileName is: %v\n", fileName)
}

func loggingFileNameSubmitted(filename string) {
	fmt.Printf("Submitted: File name: %v\n", filename)
}

func (entry *loggingFileNameEntry) validateLoggingFileName(filename string) error {
	fmt.Println("In validateLoggingFileName")
	fileName := entry.entry.Text
	if fileName == settings.JsdrSettings.Logging.LoggingFile {
		return nil
	}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("validateLoggingFileName - OpenFile error: %v", err))
		file, err = os.Create(fileName)
		if err != nil {
			sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("validateLoggingFileName - File Create error: %v", err))
			return err
		}
		file.Close()
		os.Remove(fileName)
		return nil
	}
	file.Close()
	return nil
}

// loggingLevelSelectChanged is called whenever the logging level in the logging popup is changed.
// The parameter is the new logging level.
func loggingLevelSelectChanged(level string) {
	fmt.Printf("Logging level changed to %v\n", level)
}

// closeLoggingPopUp closes the logging popup window.
func (modPopUp *modalPopUp) closeLoggingPopUp() {
	modPopUp.popUp.Hide()
}
