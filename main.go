package main

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"sway-settings/internal/config"
	"sway-settings/internal/ui"
)

func main() {
	a := app.New()
	w := a.NewWindow("Sway Settings")
	w.Resize(fyne.NewSize(620, 520))

	cfg, err := config.Load()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to load config: %w", err), w)
		cfg = &config.Config{ModKey: "Mod4", Terminal: "kitty", Opacity: 1.0}
	}

	tabs := container.NewAppTabs(
		container.NewTabItem("Variables", ui.MakeVariablesTab(cfg)),
		container.NewTabItem("Appearance", ui.MakeAppearanceTab(cfg)),
		container.NewTabItem("Outputs", ui.MakeOutputsTab(cfg)),
	)

	saveBtn := widget.NewButton("Save & Reload Sway", func() {
		if err := cfg.Save(); err != nil {
			dialog.ShowError(err, w)
			return
		}
		if err := exec.Command("swaymsg", "reload").Run(); err != nil {
			dialog.ShowError(fmt.Errorf("swaymsg reload: %w", err), w)
			return
		}
		dialog.ShowInformation("Done", "Configuration saved and Sway reloaded.", w)
	})

	w.SetContent(container.NewBorder(nil, container.NewPadded(saveBtn), nil, nil, tabs))
	w.ShowAndRun()
}
