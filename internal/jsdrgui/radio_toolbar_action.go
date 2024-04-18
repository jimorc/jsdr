package jsdrgui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type radioToolbarAction struct {
	parentWindow *fyne.Window
	action       *widget.ToolbarAction
}

// radioAction encapsulates the radio toolbar action.
var radioAction *radioToolbarAction

// newRadioToolbarAction creates a RadioToolbarAction widget.
func newRadioToolbarAction(win *fyne.Window) *radioToolbarAction {
	radioIcon := canvas.NewImageFromResource(resourceRadioSvg).Resource
	radioAction = &radioToolbarAction{parentWindow: win}
	radioAction.action = widget.NewToolbarAction(radioIcon, radioAction.radioToolbarActionActivated)
	return radioAction
}

// radioToolbarActionActivated handles mouse clicks on the radio toolbar item.
func (radioAction *radioToolbarAction) radioToolbarActionActivated() {
	// The following test is a temporary fix to disable the Radio toolbaraction until fyne issue #2306 is closed.
	if actionWin != nil {
		actionWin.window.Show()
		return
	}

	// action not disabled and no radio window exists, so create and show it.
	actionWin = newRadioWindow(radioAction.parentWindow)
	actionWin.window.Show()
}
