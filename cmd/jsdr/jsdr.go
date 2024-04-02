package main

import (
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

	// Set up program logging.
	err1 := soapylogging.CreateSoapyLogfileName(settings.JsdrSettings.Logging.LoggingFile)
	if err1 == nil {
		soapylogging.SoapyLoggingActive = true
		sdrlogger.RegisterLogHandler(soapylogging.LogSoapy)
		sdrlogger.SetLogLevel(settings.JsdrSettings.Logging.LoggingLevel)
		sdrlogger.Log(sdrlogger.Info, "jsdr Logging")
	}

	sdrApp := app.New()
	mainWindow := jsdrgui.NewMainWindow(sdrApp)
	mainWindow.Resize(fyne.NewSize(800, 300))
	mainWindow.ShowAndRun()

}
