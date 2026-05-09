package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"sway-settings/internal/config"
)

func MakeVariablesTab(cfg *config.Config) fyne.CanvasObject {
	modSelect := widget.NewSelect([]string{"Mod4 (Super)", "Mod1 (Alt)"}, func(val string) {
		if val == "Mod1 (Alt)" {
			cfg.ModKey = "Mod1"
		} else {
			cfg.ModKey = "Mod4"
		}
	})
	if cfg.ModKey == "Mod1" {
		modSelect.SetSelected("Mod1 (Alt)")
	} else {
		modSelect.SetSelected("Mod4 (Super)")
	}

	termEntry := widget.NewEntry()
	termEntry.SetText(cfg.Terminal)
	termEntry.OnChanged = func(s string) { cfg.Terminal = s }

	menuEntry := widget.NewEntry()
	menuEntry.SetText(cfg.Menu)
	menuEntry.OnChanged = func(s string) { cfg.Menu = s }

	opacityVal := widget.NewLabel(fmt.Sprintf("%.2f", cfg.Opacity))
	opacitySlider := widget.NewSlider(0.1, 1.0)
	opacitySlider.Step = 0.05
	opacitySlider.SetValue(cfg.Opacity)
	opacitySlider.OnChanged = func(v float64) {
		cfg.Opacity = v
		opacityVal.SetText(fmt.Sprintf("%.2f", v))
	}

	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Modifier key", modSelect),
			widget.NewFormItem("Terminal", termEntry),
			widget.NewFormItem("App menu command", menuEntry),
		),
		widget.NewSeparator(),
		widget.NewLabel("Default window opacity"),
		container.NewGridWithColumns(2, opacitySlider, opacityVal),
	)
}
