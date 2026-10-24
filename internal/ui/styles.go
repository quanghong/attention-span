package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Terminal layout constants tailored for 95x20 terminal window.
const (
	DefaultWidth     = 95
	DefaultHeight    = 20
	StandardBoxWidth = 91
)

// BoxWidth returns standard 91 width, gracefully adapting if terminal width is smaller.
func BoxWidth(termWidth int) int {
	if termWidth > 0 && termWidth < StandardBoxWidth+4 {
		w := termWidth - 4
		if w < 40 {
			return 40
		}
		return w
	}
	return StandardBoxWidth
}

// UI Color Palette (optimized for dark and light terminal backgrounds)
var (
	ColorFocus      = lipgloss.Color("#00D7AF") // Vibrant Mint / Emerald
	ColorFocusDark  = lipgloss.Color("#00875F")
	ColorDistract   = lipgloss.Color("#FF5F87") // Warning Coral / Crimson
	ColorDistractBg = lipgloss.Color("#5F0015")
	ColorDropMark   = lipgloss.Color("#FF8700") // Amber marker for distraction drops
	ColorElapsed    = lipgloss.Color("#5FAFFF") // Cyan / Electric Blue for elapsed target track
	ColorTrackBg    = lipgloss.Color("#3A3A3A") // Dark background track
	ColorText       = lipgloss.Color("#EEEEEE")
	ColorMuted      = lipgloss.Color("#767676")
	ColorHighlight  = lipgloss.Color("#FFFF87") // Yellow highlight
	ColorAccent     = lipgloss.Color("#AF87FF") // Lavender accent
	ColorBorder     = lipgloss.Color("#5FAFFF")
)

// Styles definitions
var (
	// Square Frame Styles
	SquareClockBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2).
			Align(lipgloss.Left)

	DistractedClockBox = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorDistract).
				Padding(1, 2).
				Align(lipgloss.Left)

	// Typography Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			MarginBottom(1)

	BigClockStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	RemainingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHighlight)

	FocusedBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorFocusDark).
			Padding(0, 1)

	DistractedBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorDistractBg).
			Padding(0, 1)

	CompletedBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5F5FAF")).
			Padding(0, 1)

	MetricLabel = lipgloss.NewStyle().
			Foreground(ColorMuted)

	MetricValue = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText)

	MetricAccent = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorFocus)

	TrackTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText)
)
