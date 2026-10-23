package ui

import (
	"fmt"
	"strings"
	"time"

	"attention-span/internal/model"
	"github.com/charmbracelet/lipgloss"
)

// RenderDurationLine renders Line 1: The total session progress from Start Time to End Time.
func RenderDurationLine(sess *model.Session, now time.Time, width int) string {
	if width < 10 {
		width = 10
	}

	progress := sess.TargetProgress(now)
	filledCount := int(progress * float64(width))
	if filledCount > width {
		filledCount = width
	}
	emptyCount := width - filledCount

	filledStyle := lipgloss.NewStyle().Foreground(ColorElapsed)
	emptyStyle := lipgloss.NewStyle().Foreground(ColorTrackBg)

	bar := filledStyle.Render(strings.Repeat("█", filledCount)) +
		emptyStyle.Render(strings.Repeat("░", emptyCount))

	elapsedStr := FormatDuration(sess.Elapsed(now))
	targetStr := FormatDuration(sess.TargetDuration)
	percentStr := fmt.Sprintf("%3.0f%%", progress*100)

	title := TrackTitle.Render("1. Target Session Line (Start → End)")
	info := fmt.Sprintf(" [%s / %s] (%s)", elapsedStr, targetStr, percentStr)

	return fmt.Sprintf("%s\n%s%s", title, bar, info)
}

// RenderAttentionRaceLine renders Line 2: The attention timeline showing focus segments
// racing with target time, and visual drop segments indicating distraction periods.
func RenderAttentionRaceLine(sess *model.Session, now time.Time, width int) string {
	if width < 10 {
		width = 10
	}

	if sess.StartTime.IsZero() || sess.TargetDuration <= 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(ColorTrackBg)
		return emptyStyle.Render(strings.Repeat("░", width))
	}

	focusStyle := lipgloss.NewStyle().Foreground(ColorFocus)
	distractStyle := lipgloss.NewStyle().Foreground(ColorDistract)
	emptyStyle := lipgloss.NewStyle().Foreground(ColorTrackBg)
	dropPointerStyle := lipgloss.NewStyle().Foreground(ColorDropMark)

	elapsed := sess.Elapsed(now)
	totalNanos := float64(sess.TargetDuration.Nanoseconds())

	// Cells buffer
	var lineBuilder strings.Builder
	markerLine := make([]rune, width)
	for i := range markerLine {
		markerLine[i] = ' '
	}

	// Build timeline cell by cell
	for i := 0; i < width; i++ {
		cellFractionStart := float64(i) / float64(width)
		cellDurationStart := time.Duration(cellFractionStart * totalNanos)

		if cellDurationStart > elapsed {
			// Future / not reached yet
			lineBuilder.WriteString(emptyStyle.Render("░"))
			continue
		}

		cellTime := sess.StartTime.Add(cellDurationStart)

		// Check if cellTime intersects with any distraction drop
		isDistracted := false
		isDropStart := false

		// Check completed drops
		for _, drop := range sess.Drops {
			if !cellTime.Before(drop.StartTime) && cellTime.Before(drop.EndTime) {
				isDistracted = true
				// Check if this cell represents the start of the drop
				if cellTime.Sub(drop.StartTime) < (sess.TargetDuration / time.Duration(width)) {
					isDropStart = true
				}
				break
			}
		}

		// Check ongoing drop
		if !isDistracted && sess.CurrentDrop != nil {
			if !cellTime.Before(sess.CurrentDrop.StartTime) {
				isDistracted = true
				if cellTime.Sub(sess.CurrentDrop.StartTime) < (sess.TargetDuration / time.Duration(width)) {
					isDropStart = true
				}
			}
		}

		if isDistracted {
			lineBuilder.WriteString(distractStyle.Render("█"))
			if isDropStart {
				markerLine[i] = '▲'
			}
		} else {
			lineBuilder.WriteString(focusStyle.Render("█"))
		}
	}

	focusTime := sess.CalculateFocusTime(now)
	focusPercent := 0.0
	if sess.TargetDuration > 0 {
		focusPercent = (float64(focusTime) / float64(sess.TargetDuration)) * 100
	}
	if focusPercent > 100 {
		focusPercent = 100
	}

	title := TrackTitle.Render("2. Attention Race & Distraction Drops")
	info := fmt.Sprintf(" [Focus: %s] (%3.0f%% of target)", FormatDuration(focusTime), focusPercent)

	line2 := lineBuilder.String() + info

	// Drop markers indicator line
	hasMarkers := false
	for _, r := range markerLine {
		if r != ' ' {
			hasMarkers = true
			break
		}
	}

	totalDrops := len(sess.Drops)
	if sess.CurrentDrop != nil {
		totalDrops++
	}

	var markerRow string
	if hasMarkers || totalDrops > 0 {
		dropsSummary := fmt.Sprintf(" (%d drops · %s distracted)", totalDrops, FormatDuration(sess.TotalDistractionTime(now)))
		markerRow = "\n" + dropPointerStyle.Render(string(markerLine)) + MetricLabel.Render(dropsSummary)
	}

	return fmt.Sprintf("%s\n%s%s", title, line2, markerRow)
}

// FormatDuration formats duration nicely into "01h 30m 00s" or "45m 12s".
func FormatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%02dh %02dm %02ds", h, m, s)
	}
	return fmt.Sprintf("%02dm %02ds", m, s)
}

// FormatDurationShort formats into compact "1h 30m" or "45m".
func FormatDurationShort(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	h := d / time.Hour
	m := (d % time.Hour) / time.Minute
	if h > 0 && m > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	if h > 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dm", m)
}
