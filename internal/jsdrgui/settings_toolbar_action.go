package jsdrgui

import (
	"fmt"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// NewSettingsToolbarAction creates a SettingsToolbarAction widget
func NewSettingsToolbarAction() *widget.ToolbarAction {
	settingsIcon := canvas.NewImageFromResource(resourceSettingsSvg).Resource
	return widget.NewToolbarAction(settingsIcon, settingsToolbarActionActivated)
}

func settingsToolbarActionActivated() {
	fmt.Println("In settingsToolbarActionActivated")
}
