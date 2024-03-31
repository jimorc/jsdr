// Package jsdrgui contains the detailed widgets used in the go_sdr app.
package jsdrgui

import (
	"fmt"
	"internal/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// StartStopToolbarAction is a gui.TwoStateToolbarAction that defines the actions performed by the SDR start/stop toolbar button.
type StartStopToolbarAction struct {
	action *gui.TwoStateToolbarAction
}

// NewStartStopToolbarAction creates a StartStopToolbarAction object.
func NewStartStopToolbarAction() *StartStopToolbarAction {
	startIcon := canvas.NewImageFromResource(resourceStartSvg).Resource

	stopIcon := canvas.NewImageFromResource(resourceStopSvg).Resource

	startStop := StartStopToolbarAction{}
	startStop.action = gui.NewTwoStateToolbarAction(startIcon, stopIcon, startStop.startActivated, startStop.stopActivated)
	return &startStop
}

// ToolbarObject returns a pointer to the underlying TwoStateToolbarObject
func (t *StartStopToolbarAction) ToolbarObject() fyne.CanvasObject {
	return t.action.ToolbarObject()
}

func (t *StartStopToolbarAction) startActivated() {
	fmt.Println("In startActivated")
}

func (t *StartStopToolbarAction) stopActivated() {
	fmt.Println("In stopActivated")
}
