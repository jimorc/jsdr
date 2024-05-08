package jsdrgui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type sourceToolbarAction struct {
	parentWindow *fyne.Window
	action       *widget.ToolbarAction
}

// newSourceToolbarAction creates a s ource toolbar action widget.
func newSourceToolbarAction(win *fyne.Window) *sourceToolbarAction {
	sourceIcon := canvas.NewImageFromResource(resourceBlackWrenchSvg).Resource
	sourceAction := &sourceToolbarAction{parentWindow: win}
	sourceAction.action = widget.NewToolbarAction(sourceIcon, sourceAction.sourceToolbarActionActivated)
	return sourceAction
}

// sourceToolbarActionActivated handles mouse clicks on the source toolbar item.
func (sourceAction *sourceToolbarAction) sourceToolbarActionActivated() {
	// The following test is a temporary fix to disable the Hardware toolbaraction until fyne issue #2306 is closed.
	if actionWin != nil {
		actionWin.window.Show()
		return
	}

	// action not disabled and no hardware window exists, so create and show it.
	actionWin = newSourceWindow(sourceAction.parentWindow)
	if actionWin != nil {
		actionWin.window.Show()
	}
}
