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
	if width < 20 {
		width = 20
	}

	progress := sess.TargetProgress(now)
	elapsedStr := FormatDuration(sess.Elapsed(now))
	targetStr := FormatDuration(sess.TargetDuration)
	percentStr := fmt.Sprintf("%3.0f%%", progress*100)

	title := TrackTitle.Render("1. Target Session Line (Start → End)")
	info := fmt.Sprintf(" [%s / %s] (%s)", elapsedStr, targetStr, percentStr)
	infoWidth := lipgloss.Width(info)

	barLen := width - infoWidth
	if barLen < 10 {
		barLen = 10
	}
	filledCount := int(progress * float64(barLen))
	if filledCount > barLen {
		filledCount = barLen
	}
	emptyCount := barLen - filledCount

	filledStyle := lipgloss.NewStyle().Foreground(ColorElapsed)
	emptyStyle := lipgloss.NewStyle().Foreground(ColorTrackBg)

	bar := filledStyle.Render(strings.Repeat("█", filledCount)) +
		emptyStyle.Render(strings.Repeat("░", emptyCount))

	return fmt.Sprintf("%s\n%s%s", title, bar, info)
}

// RenderAttentionRaceLine renders Line 2: The attention timeline showing focus segments
// racing with target time, and visual drop segments indicating distraction periods.
func RenderAttentionRaceLine(sess *model.Session, now time.Time, width int) string {
	if width < 20 {
		width = 20
	}

	focusTime := sess.CalculateFocusTime(now)
	focusPercent := sess.FocusProgress(now) * 100

	title := TrackTitle.Render("2. Attention Race & Distraction Drops")
	info := fmt.Sprintf(" [Focus: %s] (%3.0f%% of target)", FormatDuration(focusTime), focusPercent)
	infoWidth := lipgloss.Width(info)

	barLen := width - infoWidth
	if barLen < 10 {
		barLen = 10
	}

	emptyStyle := lipgloss.NewStyle().Foreground(ColorTrackBg)
	if sess.StartTime.IsZero() || sess.TargetDuration <= 0 {
		bar := emptyStyle.Render(strings.Repeat("░", barLen))
		return fmt.Sprintf("%s\n%s%s", title, bar, info)
	}

	focusStyle := lipgloss.NewStyle().Foreground(ColorFocus)
	distractStyle := lipgloss.NewStyle().Foreground(ColorDistract)
	dropPointerStyle := lipgloss.NewStyle().Foreground(ColorDropMark)

	elapsed := sess.Elapsed(now)
	totalNanos := float64(sess.TargetDuration.Nanoseconds())

	var lineBuilder strings.Builder
	markerLine := make([]rune, barLen)
	for i := range markerLine {
		markerLine[i] = ' '
	}

	// Build timeline cell by cell
	for i := 0; i < barLen; i++ {
		cellFractionStart := float64(i) / float64(barLen)
		cellDurationStart := time.Duration(cellFractionStart * totalNanos)

		if cellDurationStart > elapsed {
			lineBuilder.WriteString(emptyStyle.Render("░"))
			continue
		}

		cellTime := sess.StartTime.Add(cellDurationStart)

		isDistracted := false
		isDropStart := false

		for _, drop := range sess.Drops {
			if !cellTime.Before(drop.StartTime) && cellTime.Before(drop.EndTime) {
				isDistracted = true
				if cellTime.Sub(drop.StartTime) < (sess.TargetDuration / time.Duration(barLen)) {
					isDropStart = true
				}
				break
			}
		}

		if !isDistracted && sess.CurrentDrop != nil {
			if !cellTime.Before(sess.CurrentDrop.StartTime) {
				isDistracted = true
				if cellTime.Sub(sess.CurrentDrop.StartTime) < (sess.TargetDuration / time.Duration(barLen)) {
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

	totalDrops := len(sess.Drops)
	if sess.CurrentDrop != nil {
		totalDrops++
	}

	hasMarkers := false
	for _, r := range markerLine {
		if r != ' ' {
			hasMarkers = true
			break
		}
	}

	var markerRow string
	if hasMarkers || totalDrops > 0 {
		dropsSummary := fmt.Sprintf(" (%d drops · %s distracted)", totalDrops, FormatDuration(sess.TotalDistractionTime(now)))
		markerRow = "\n" + dropPointerStyle.Render(string(markerLine)) + MetricLabel.Render(dropsSummary)
	}

	return fmt.Sprintf("%s\n%s%s%s", title, lineBuilder.String(), info, markerRow)
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
