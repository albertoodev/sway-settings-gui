package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"sway-settings/internal/config"
)

func MakeOutputsTab(cfg *config.Config) fyne.CanvasObject {
	if len(cfg.Outputs) == 0 {
		return container.NewCenter(widget.NewLabel("No outputs found in outputs.conf"))
	}

	transforms := []string{"", "90", "180", "270", "flipped", "flipped-90", "flipped-180", "flipped-270"}

	var rows []fyne.CanvasObject
	for i := range cfg.Outputs {
		out := &cfg.Outputs[i]

		nameLabel := widget.NewLabelWithStyle(out.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		resEntry := widget.NewEntry()
		resEntry.SetText(out.Resolution)
		resEntry.PlaceHolder = "1920x1080"
		resEntry.OnChanged = func(s string) { out.Resolution = s }

		posEntry := widget.NewEntry()
		posEntry.SetText(out.Position)
		posEntry.PlaceHolder = "0,0"
		posEntry.OnChanged = func(s string) { out.Position = s }

		transformSel := widget.NewSelect(transforms, func(s string) { out.Transform = s })
		transformSel.SetSelected(out.Transform)

		rows = append(rows,
			nameLabel,
			widget.NewForm(
				widget.NewFormItem("Resolution", resEntry),
				widget.NewFormItem("Position", posEntry),
				widget.NewFormItem("Transform", transformSel),
			),
			widget.NewSeparator(),
		)
	}

	return container.NewVScroll(container.NewVBox(rows...))
}
