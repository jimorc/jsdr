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

// NewLoggingToolbarAction creates a SettingsToolbarAction widget.
func NewLoggingToolbarAction(win *fyne.Window) *widget.ToolbarAction {
	loggingIcon := canvas.NewImageFromResource(resourceLoggingSvg).Resource
	loggingAction := loggingToolbarAction{parentWindow: win}
	loggingAction.action = widget.NewToolbarAction(loggingIcon, loggingAction.loggingToolbarActionActivated)
	return loggingAction.action
}

func (loggingAction *loggingToolbarAction) loggingToolbarActionActivated() {

	loggingPopup := newSDRLoggerSettingsPopUp(loggingAction.parentWindow)
	loggingPopup.Show()

	fmt.Println("In settingsToolbarActionActivated")
}

func (loggingAction *loggingToolbarAction) loggingConfirmCallback() {
	fmt.Println("In loggingConfirmCallback")
}
