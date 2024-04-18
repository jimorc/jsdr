package jsdrgui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type radioToolbarAction struct {
	parentWindow *fyne.Window
	action       *widget.ToolbarAction
	disabled     bool
}

// RadioAction encapsulates the radio toolbar action.
var RadioAction radioToolbarAction

// newRadioToolbarAction creates a RadioToolbarAction widget.
func newRadioToolbarAction(win *fyne.Window) *radioToolbarAction {
	radioIcon := canvas.NewImageFromResource(resourceRadioSvg).Resource
	RadioAction = radioToolbarAction{parentWindow: win, disabled: false}
	RadioAction.action = widget.NewToolbarAction(radioIcon, RadioAction.radioToolbarActionActivated)
	return &RadioAction
}

// radioToolbarActionActivated handles mouse clicks on the radio toolbar item.
func (radioAction *radioToolbarAction) radioToolbarActionActivated() {
	// The following test is a temporary fix to disable the Radio toolbaraction until fyne issue #2306 is closed.
	if RadioAction.disabled {
		if RadioWin != nil {
			RadioWin.Window.Show()
		}
		return
	}

	// action not disabled and no radio window exists, so create and show it.
	RadioWin = newRadioWindow(radioAction.parentWindow)
	RadioWin.Window.Show()
	disableMainToolbar()
}

// disable disables the radio toolbar action. This is used to prevent diaplaying multiple windows on top of the main window.
func (radioAction *radioToolbarAction) disable() {
	RadioAction.disabled = true
}

// enable enables the radio toolbar action.
func (radioAction *radioToolbarAction) enable() {
	RadioAction.disabled = false
}
