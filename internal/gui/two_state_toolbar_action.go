// Package gui provides the general widgets, that are not defined in fyne, that are used in go_sdr app.
//
// For example, the TwoStateToolbarAction widget provides a general ToolbarAction that can have two states.
// This differs from the object created by this widget which provides the functionality specific to that object.
// That object is defined in the package "gosdrgui"
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TwoStateToolbarAction is a two state version of the ToolbarAction.
// It displays one of two different icons and calls one of two functions based on its state.
type TwoStateToolbarAction struct {
	state         bool
	uncheckedIcon fyne.Resource
	checkedIcon   fyne.Resource
	button        *widget.Button

	onUncheckedActivated func()
	onCheckedActivated   func()
	onActivated          func()
}

// NewTwoStateToolbarAction creates a TwoStateToolbarAction object
func NewTwoStateToolbarAction(unchecked fyne.Resource, checked fyne.Resource,
	uncheckedActivated func(), checkedActivated func()) *TwoStateToolbarAction {
	w := &TwoStateToolbarAction{
		state:                false,
		uncheckedIcon:        unchecked,
		checkedIcon:          checked,
		onUncheckedActivated: uncheckedActivated,
		onCheckedActivated:   checkedActivated,
	}
	w.onActivated = w.activated
	w.button = widget.NewButtonWithIcon("", w.uncheckedIcon, w.onActivated)

	return w
}

// ToolbarObject sets the icon to be displayed by the object based on the object's state and returns the button representing
// the object.
func (t *TwoStateToolbarAction) ToolbarObject() fyne.CanvasObject {
	if t.state {
		t.button.SetIcon(t.checkedIcon)
	} else {
		t.button.SetIcon(t.uncheckedIcon)
	}
	return t.button
}

// onActivated is executed whenever the the action button is clicked. It does the following:
// 1. Toggles the state of the action.
// 2. Calls the action based on its previous state.
// 3. Refreshes the display of the action object to reflect its current state.
func (t *TwoStateToolbarAction) activated() {
	if t.state {
		t.state = false
		t.onCheckedActivated()
	} else {
		t.state = true
		t.onUncheckedActivated()
	}
	t.ToolbarObject().Refresh()
}
