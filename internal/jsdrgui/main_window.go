package jsdrgui

import (
	"fmt"
	"internal/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

// NewMainWindow creates the main window for the go_sdr app.
func NewMainWindow(sdrApp fyne.App) fyne.Window {
	mainWin := sdrApp.NewWindow("jsdr")
	startStop := NewStartStopToolbarAction()
	loggingToolbarAction := NewLoggingToolbarAction(&mainWin)
	radioToolbarAction := NewRadioToolbarAction(&mainWin)
	toolBar := widget.NewToolbar(radioToolbarAction, startStop, loggingToolbarAction)

	mainWin.SetContent(toolBar)
	mainWin.SetOnClosed(mainWindowClosing)
	return mainWin
}

// mainWindowClosing
func mainWindowClosing() {
	fmt.Println("In mainWindowClosing")
	if RadioWin != nil {
		RadioWin.Window.Close()
	}
	err := settings.JsdrSettings.Save()
	if err == nil {
		sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Settings saved to %v", settings.SettingsFileName))
	} else {
		sdrlogger.Log(sdrlogger.Error, fmt.Sprintf("Unable to save settings file: %v\n    %v", settings.SettingsFileName, err))
	}
}
