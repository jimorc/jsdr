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

// RadioAction encapsulates the radio toolbar action.
var RadioAction = radioToolbarAction{}

// NewRadioToolbarAction creates a RadioToolbarAction widget.
func NewRadioToolbarAction(win *fyne.Window) *widget.ToolbarAction {
	radioIcon := canvas.NewImageFromResource(resourceRadioSvg).Resource
	RadioAction := radioToolbarAction{parentWindow: win}
	RadioAction.action = widget.NewToolbarAction(radioIcon, RadioAction.radioToolbarActionActivated)
	return RadioAction.action
}

func (radioAction *radioToolbarAction) radioToolbarActionActivated() {
	// The following test is a temporary fix to disable the Radio toolbaraction until fyne issue #2306 is closed.
	if RadioWin != nil {
		return
	}

	RadioWin = newRadioWindow(radioAction.parentWindow)
	RadioWin.Window.Show()
}
