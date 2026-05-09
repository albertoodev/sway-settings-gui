# Sway Settings

A lightweight, Go-based GUI configuration tool built specifically for my sway config in [dotfiles](https://github.com/albertoodev/.dotfiles). It manages [Sway](https://swaywm.org/) and [SwayFX](https://github.com/WillPower3309/swayfx) settings using the [Fyne](https://fyne.io/) toolkit.

## Features

- **Variables Management**: Configure your modifier key, preferred terminal, application menu command, and default window opacity.
- **Appearance (SwayFX Support)**: 
  - Adjust inner and outer gaps.
  - Set corner radius for a rounded look.
  - Configure border widths.
  - Enable and blur effects (passes and radius).
  - Adjust dimming for inactive windows.
- **Output Configuration**: Manage multiple displays, including resolution, position, and transforms (rotation/flipping).
- **One-Click Apply**: Save settings and reload Sway instantly via `swaymsg`.

## Prerequisites

- **Sway** or **SwayFX** (required for `swaymsg` and specific appearance features).
- **Go** 1.21+ (to build from source).
- **Fyne Dependencies**: See the [Fyne setup guide](https://developer.fyne.io/started/) for your platform (usually requires development headers for X11/Wayland, GL, etc.).

## Configuration Structure

The tool manages configuration across three main files in `~/.config/sway/`:

1.  **`config`**: Handles core variables (`$mod`, `$term`, `$menu`, `$opacity`).
2.  **`config.d/appearance.conf`**: Manages aesthetic settings (gaps, borders, blur, dimming).
3.  **`config.d/outputs.conf`**: Manages display-specific settings.

