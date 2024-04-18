package jsdrgui

import (
	"fmt"
	"internal/settings"

	"fyne.io/fyne/v2"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

var mainWin fyne.Window = nil

// NewMainWindow creates the main window for the go_sdr app.
func NewMainWindow(sdrApp fyne.App) fyne.Window {
	mainWin = sdrApp.NewWindow("jsdr")
	mainToolBar := newMainToolbar(&mainWin)
	mainWin.SetContent(mainToolBar)
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
