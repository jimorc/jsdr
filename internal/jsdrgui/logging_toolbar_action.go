package jsdrgui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type loggingToolbarAction struct {
	parentWindow *fyne.Window
	action       *widget.ToolbarAction
	disabled     bool
}

// radioAction encapsulates the radio toolbar action.
var loggingAction *loggingToolbarAction

// newLoggingToolbarAction creates a SettingsToolbarAction widget.
func newLoggingToolbarAction(win *fyne.Window) *loggingToolbarAction {
	loggingIcon := canvas.NewImageFromResource(resourceLogsettingsSvg).Resource
	loggingAction = &loggingToolbarAction{parentWindow: win, disabled: false}
	loggingAction.action = widget.NewToolbarAction(loggingIcon, loggingAction.loggingToolbarActionActivated)
	return loggingAction
}

func (loggingAction *loggingToolbarAction) loggingToolbarActionActivated() {
	if loggingAction.disabled {
		if loggingWin != nil {
			loggingWin.window.Show()
		}
		return
	}

	loggingWin := newSDRLoggerSettingsWindow()
	loggingWin.window.Show()
	disableMainToolbar()

	fmt.Println("In settingsToolbarActionActivated")
}

// disable disables the logging toolbar action. This is used to prevent displaying multiple windows on top of the main window.
func (loggingAction *loggingToolbarAction) disable() {
	loggingAction.disabled = true
}

// enable enables the logging toolbar action.
func (loggingAction *loggingToolbarAction) enable() {
	loggingAction.disabled = false
}
