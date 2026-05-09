package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"sway-settings/internal/config"
)

func MakeAppearanceTab(cfg *config.Config) fyne.CanvasObject {
	intRow := func(label string, min, max float64, get func() int, set func(int)) fyne.CanvasObject {
		lbl := widget.NewLabel(fmt.Sprintf("%d", get()))
		s := widget.NewSlider(min, max)
		s.Step = 1
		s.SetValue(float64(get()))
		s.OnChanged = func(v float64) {
			set(int(v))
			lbl.SetText(fmt.Sprintf("%d", int(v)))
		}
		return container.NewGridWithColumns(3, widget.NewLabel(label), s, lbl)
	}

	floatRow := func(label string, min, max, step float64, get func() float64, set func(float64)) fyne.CanvasObject {
		lbl := widget.NewLabel(fmt.Sprintf("%.2f", get()))
		s := widget.NewSlider(min, max)
		s.Step = step
		s.SetValue(get())
		s.OnChanged = func(v float64) {
			set(v)
			lbl.SetText(fmt.Sprintf("%.2f", v))
		}
		return container.NewGridWithColumns(3, widget.NewLabel(label), s, lbl)
	}

	blurCheck := widget.NewCheck("Enable blur (SwayFX)", func(v bool) { cfg.BlurEnabled = v })
	blurCheck.SetChecked(cfg.BlurEnabled)

	return container.NewVBox(
		widget.NewLabelWithStyle("Gaps & Borders", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		intRow("Gaps inner (px)", 0, 50, func() int { return cfg.GapsInner }, func(v int) { cfg.GapsInner = v }),
		intRow("Gaps outer (px)", 0, 50, func() int { return cfg.GapsOuter }, func(v int) { cfg.GapsOuter = v }),
		intRow("Corner radius (px)", 0, 30, func() int { return cfg.CornerRadius }, func(v int) { cfg.CornerRadius = v }),
		intRow("Border width (px)", 0, 5, func() int { return cfg.BorderSize }, func(v int) { cfg.BorderSize = v }),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Blur (SwayFX)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		blurCheck,
		intRow("Blur passes", 1, 10, func() int { return cfg.BlurPasses }, func(v int) { cfg.BlurPasses = v }),
		intRow("Blur radius", 1, 20, func() int { return cfg.BlurRadius }, func(v int) { cfg.BlurRadius = v }),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Inactive Windows", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		floatRow("Dim inactive", 0, 0.5, 0.05, func() float64 { return cfg.DimInactive }, func(v float64) { cfg.DimInactive = v }),
	)
}
