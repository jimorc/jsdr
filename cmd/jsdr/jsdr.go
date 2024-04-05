package main

import (
	"fmt"
	"internal/jsdrgui"
	"internal/settings"
	"internal/soapylogging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/pothosware/go-soapy-sdr/pkg/sdrlogger"
)

func main() {
	// Create and load program settings.
	settings.JsdrSettings = settings.NewSettings()
	err := settings.JsdrSettings.Load()

	// Set up program logging.
	soapylogging.SoapyLoggingActive = true
	sdrlogger.RegisterLogHandler(soapylogging.LogSoapy)
	sdrlogger.SetLogLevel(settings.JsdrSettings.Logging.LoggingLevel)
	sdrlogger.Log(sdrlogger.Info, "jsdr Logging initialized")
	if err != nil {
		sdrlogger.Log(sdrlogger.Error, fmt.Sprintf("Unable to load settings file:\n    %v", err))

	// Create and show program GUI
	sdrApp := app.New()
	mainWindow := jsdrgui.NewMainWindow(sdrApp)
	mainWindow.Resize(fyne.NewSize(800, 300))
	mainWindow.ShowAndRun()

}
