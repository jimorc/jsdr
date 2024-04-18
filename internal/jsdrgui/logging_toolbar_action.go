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
}

// loggingAction encapsulates the logging toolbar action.
var loggingAction *loggingToolbarAction

// newLoggingToolbarAction creates a SettingsToolbarAction widget.
func newLoggingToolbarAction(win *fyne.Window) *loggingToolbarAction {
	loggingIcon := canvas.NewImageFromResource(resourceLogsettingsSvg).Resource
	loggingAction = &loggingToolbarAction{parentWindow: win}
	loggingAction.action = widget.NewToolbarAction(loggingIcon, loggingAction.loggingToolbarActionActivated)
	return loggingAction
}

func (loggingAction *loggingToolbarAction) loggingToolbarActionActivated() {
	if actionWin != nil {
		actionWin.window.Show()
		return
	}

	actionWin = newSDRLoggerSettingsWindow()
	actionWin.window.Show()

	fmt.Println("In settingsToolbarActionActivated")
}
