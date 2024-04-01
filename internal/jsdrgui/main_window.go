package jsdrgui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// NewMainWindow creates the main window for the go_sdr app.
func NewMainWindow(sdrApp fyne.App) fyne.Window {
	mainWin := sdrApp.NewWindow("jsdr")
	startStop := NewStartStopToolbarAction()
	loggingToolbarAction := NewLoggingToolbarAction(&mainWin)
	toolBar := widget.NewToolbar(loggingToolbarAction, startStop)

	mainWin.SetContent(toolBar)

	return mainWin
}
