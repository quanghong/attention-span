package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// RenderSetupView displays the timer configuration screen before starting.
func RenderSetupView(target time.Duration, width, height int) string {
	boxWidth := 56
	if width > 0 && width < 60 {
		boxWidth = width - 4
	}

	title := TitleStyle.Render("⚡ ATTENTION SPAN TRACKER ⚡")
	subtitle := MetricLabel.Render("Set your work target duration and train your focus")

	hours := target / time.Hour
	minutes := (target % time.Hour) / time.Minute

	durationDisplay := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorHighlight).
		Background(lipgloss.Color("#262626")).
		Padding(0, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Render(fmt.Sprintf("%02dh %02dm", hours, minutes))

	adjustPrompt := lipgloss.NewStyle().
		Foreground(ColorText).
		Render("Use Arrow keys or H/J/K/L to adjust target time")

	quickPresets := MetricLabel.Render(
		"Quick Presets: [1] 25m  [2] 50m  [3] 1h  [4] 1.5h  [5] 2h  [6] 3h  [7] 4h",
	)

	startNotice := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorFocus).
		Render("Press [ENTER] or [SPACE] to START session")

	controls := MetricLabel.Render(
		"[← / →] ±15m   [↑ / ↓] ±1h   [h] History   [q] Quit",
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		"",
		durationDisplay,
		"",
		adjustPrompt,
		quickPresets,
		"",
		startNotice,
		"",
		controls,
	)

	boxStyle := SquareClockBox.Copy().Width(boxWidth)
	card := boxStyle.Render(content)

	if height > 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
	}

	return card
}
