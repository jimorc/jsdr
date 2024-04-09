package main

import (
	"fmt"
	"internal/jsdrgui"
	"internal/settings"
	"internal/soapylogging"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

func main() {
	// Create and load program settings.
	settings.JsdrSettings = settings.NewSettings()
	err := settings.JsdrSettings.Load()

	// Set up program logging.
	soapylogging.CreateSoapyLogFile()
	sdrlogger.RegisterLogHandler(soapylogging.LogSoapy)
	sdrlogger.SetLogLevel(sdrlogger.SDRLogLevel(atomic.LoadInt64(&settings.JsdrSettings.LoggingLevel)))
	sdrlogger.Log(sdrlogger.Info, "jsdr Logging initialized")
	if err == nil {
		sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Settings loaded from %v", settings.SettingsFileName))
	} else {
		sdrlogger.Log(sdrlogger.Error, fmt.Sprintf("Unable to load settings file: %v\n    %v", settings.SettingsFileName, err))
	}
	// Create and show program GUI
	sdrApp := app.New()
	mainWindow := jsdrgui.NewMainWindow(sdrApp)
	mainWindow.Resize(fyne.NewSize(800, 300))
	mainWindow.SetOnClosed(mainWindowClosing)
	mainWindow.ShowAndRun()
}

// mainWindowClosing
func mainWindowClosing() {
	fmt.Println("In mainWindowClosing")
	err := settings.JsdrSettings.Save()
	if err == nil {
		sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Settings saved to %v", settings.SettingsFileName))
	} else {
		sdrlogger.Log(sdrlogger.Error, fmt.Sprintf("Unable to save settings file: %v\n    %v", settings.SettingsFileName, err))
	}
}
