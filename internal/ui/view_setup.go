package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// RenderSetupView displays the timer configuration screen before starting.
func RenderSetupView(target time.Duration, width, height int) string {
	boxWidth := BoxWidth(width)

	title := TitleStyle.Render("ATTENTION SPAN TRACKER")
	subtitle := MetricLabel.Render("Set your work target duration and train your focus")

	hours := target / time.Hour
	minutes := (target % time.Hour) / time.Minute

	durationBadge := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorHighlight).
		Background(lipgloss.Color("#262626")).
		Padding(0, 1).
		Render(fmt.Sprintf("%02dh %02dm", hours, minutes))

	durationLine := MetricLabel.Render("Target Duration: ") + durationBadge

	startNotice := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorFocus).
		Render("▶ Press [ENTER] or [SPACE] to START session")

	controls := MetricLabel.Render(
		"Controls: [←/→] ±15m   [↑/↓] ±1h   [h] History   [q] Quit",
	)

	elements := []string{
		title,
		subtitle,
		"",
		durationLine,
		"",
		startNotice,
		"",
		controls,
	}

	content := lipgloss.JoinVertical(lipgloss.Left, elements...)
	card := SquareClockBox.Width(boxWidth).Render(content)

	w, h := width, height
	if w <= 0 {
		w = DefaultWidth
	}
	if h <= 0 {
		h = DefaultHeight
	}

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, card)
}
