package main

import (
	"fmt"
	"internal/jsdrgui"
	"internal/settings"
	"internal/soapylogging"
	"log"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

func main() {
	// Create and load program settings.
	settings.JsdrSettings = settings.NewSettings()
	err := settings.JsdrSettings.Load()

	// Set up program logging.
	// If error, then logging cannot be done, so we terminate with error.
	err1 := soapylogging.CreateSoapyLogFile()
	if err1 != nil {
		log.Fatalln("Unable to create the logging file.")
	}
	sdrlogger.RegisterLogHandler(soapylogging.LogSoapy)
	sdrlogger.SetLogLevel(sdrlogger.SDRLogLevel(atomic.LoadInt64(&settings.JsdrSettings.LoggingLevel)))
	sdrlogger.Log(sdrlogger.Info, "jsdr Logging initialized")
	if err == nil {
		sdrlogger.Log(sdrlogger.Trace, fmt.Sprintf("Settings loaded from %v", settings.SettingsFileName))
	} else {
		sdrlogger.Log(sdrlogger.Error, fmt.Sprintf("Unable to load settings file: %v\n    %v", settings.SettingsFileName, err))
	}

	jsdrgui.SdrApp.Settings().SetTheme(theme.LightTheme())
	// Create and show program GUI
	mainWindow := jsdrgui.NewMainWindow(jsdrgui.SdrApp)
	mainWindow.Resize(fyne.NewSize(800, 300))
	mainWindow.ShowAndRun()
}
