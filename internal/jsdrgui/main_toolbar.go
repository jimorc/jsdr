package jsdrgui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type actionWindow struct {
	window fyne.Window
}

// actionWin contains the window that is created by any of the toolbar actions.
// This variable is used to ensure that only one action window is displayed at a time.
// This is used to disable the toolbar actions.
// TODO: When fyne supports disabling the toolbar or toolbar actions, this code could be replaced.
var actionWin *actionWindow = nil

var mainToolbar *widget.Toolbar = nil
var startStop = NewStartStopToolbarAction()

// newMainToolbar creates the main toolbar
func newMainToolbar(mainWin *fyne.Window) *widget.Toolbar {
	loggingAction := newLoggingToolbarAction(mainWin)
	radioAction := newRadioToolbarAction(mainWin)
	mainToolbar = widget.NewToolbar(radioAction.action, startStop, loggingAction.action)
	return mainToolbar
}
