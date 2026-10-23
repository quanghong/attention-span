package ui

import (
	"github.com/charmbracelet/lipgloss"
)

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
			Align(lipgloss.Center)

	DistractedClockBox = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorDistract).
				Padding(1, 2).
				Align(lipgloss.Center)

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

	PausedBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorHighlight).
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

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	TrackTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText)
)
