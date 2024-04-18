package jsdrgui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type actionWindow struct {
	window fyne.Window
}

var actionWin *actionWindow = nil

var mainToolbar *widget.Toolbar = nil
var startStop = NewStartStopToolbarAction()

// newMainToolbar creates the main toolbar
func newMainToolbar(mainWin *fyne.Window) *widget.Toolbar {
	loggingAction = newLoggingToolbarAction(mainWin)
	radioAction = newRadioToolbarAction(mainWin)
	mainToolbar = widget.NewToolbar(radioAction.action, startStop, loggingAction.action)
	return mainToolbar
}

// disableMainToolbar disables each of the toolbar items in the toolbar.
// This function should be called whenever a window is displayed over the main window.
func disableMainToolbar() {
	radioAction.disable()
	loggingAction.disable()
}

// enableMainToolbar enables each of the toolbar items in the toolbar.
// This function should be called whenever a window that is displayed over the main window is closed.
func enableMainToolbar() {
	radioAction.enable()
	loggingAction.enable()
}
