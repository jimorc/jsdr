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

// NewRadioToolbarAction creates a RadioToolbarAction widget.
func NewRadioToolbarAction(win *fyne.Window) *widget.ToolbarAction {
	radioIcon := canvas.NewImageFromResource(resourceRadioSvg).Resource
	radioAction := radioToolbarAction{parentWindow: win}
	radioAction.action = widget.NewToolbarAction(radioIcon, radioAction.radioToolbarActionActivated)
	return radioAction.action
}

func (radioAction *radioToolbarAction) radioToolbarActionActivated() {

	radioPopup := newRadioPopUp(radioAction.parentWindow)
	radioPopup.Show()
}
