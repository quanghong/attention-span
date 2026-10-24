package ui

import (
	"fmt"
	"time"

	"attention-span/internal/model"

	"github.com/charmbracelet/lipgloss"
)

// RenderTimerView renders the primary active session screen with the Square Clock
// and the two racing lines (Target Session Line and Attention Line with Distraction Drops).
func RenderTimerView(sess *model.Session, now time.Time, termWidth, termHeight int) string {
	boxWidth := BoxWidth(termWidth)
	trackWidth := boxWidth - 6
	if trackWidth < 20 {
		trackWidth = 20
	}

	currentTimeStr := now.Format("03:04:05 PM")
	targetEndTime := sess.StartTime.Add(sess.TargetDuration)
	endTimeStr := targetEndTime.Format("03:04 PM")
	remaining := sess.Remaining(now)

	clockHeader := BigClockStyle.Render("🕒 "+currentTimeStr) +
		MetricLabel.Render("   │   Target End: ") + MetricValue.Render(endTimeStr) +
		MetricLabel.Render("   │   ") + RemainingStyle.Render("Remaining: "+FormatDuration(remaining))

	var badge string
	var statusMsg string
	if sess.State == model.StateFocused {
		badge = FocusedBadge.Render("● IN THE ZONE (FOCUSED)")
		statusMsg = lipgloss.NewStyle().Foreground(ColorFocus).Render("Press [SPACE] or [d] when you get distracted")
	} else if sess.State == model.StateDistracted {
		badge = DistractedBadge.Render("▲ DISTRACTED (ATTENTION ON HOLD)")
		currentDropDuration := ""
		if sess.CurrentDrop != nil {
			currentDropDuration = fmt.Sprintf(" (%s elapsed)", FormatDuration(sess.CurrentDrop.CalculateDuration(now)))
		}
		statusMsg = lipgloss.NewStyle().Foreground(ColorDistract).Render("Press [SPACE] or [d] when you refocus" + currentDropDuration)
	} else if sess.State == model.StateCompleted {
		badge = CompletedBadge.Render("★ TARGET TIME REACHED")
		statusMsg = MetricAccent.Render("Session duration reached! Press [e] to view summary.")
	}

	statusRow := badge + "  " + statusMsg

	// Line 1: Duration from Start Time to End Time
	line1 := RenderDurationLine(sess, now, trackWidth)

	// Line 2: Attention Race with Distraction Drops
	line2 := RenderAttentionRaceLine(sess, now, trackWidth)

	// Metrics Row
	focusRatio := sess.CalculateFocusRatio(now) * 100
	totalFocus := sess.CalculateFocusTime(now)
	totalDistract := sess.TotalDistractionTime(now)
	dropCount := len(sess.Drops)
	if sess.CurrentDrop != nil {
		dropCount++
	}

	metricsBar := fmt.Sprintf(
		"%s %s   %s %s   %s %s   %s %s",
		MetricLabel.Render("Ratio:"),
		MetricAccent.Render(fmt.Sprintf("%4.1f%%", focusRatio)),
		MetricLabel.Render("Focused:"),
		MetricValue.Render(FormatDurationShort(totalFocus)),
		MetricLabel.Render("Distracted:"),
		lipgloss.NewStyle().Foreground(ColorDistract).Render(FormatDurationShort(totalDistract)),
		MetricLabel.Render("Drops:"),
		MetricValue.Render(fmt.Sprintf("%d", dropCount)),
	)

	// Controls Footer
	controls := MetricLabel.Render(
		"[Space] Toggle Attention   │   [e] End Early   │   [q] Quit",
	)

	contentElements := []string{
		clockHeader,
		"",
		statusRow,
		"",
		line1,
		"",
		line2,
		"",
		metricsBar,
		"",
		controls,
	}

	content := lipgloss.JoinVertical(lipgloss.Left, contentElements...)

	var boxStyle lipgloss.Style
	if sess.State == model.StateDistracted {
		boxStyle = DistractedClockBox.Width(boxWidth)
	} else {
		boxStyle = SquareClockBox.Width(boxWidth)
	}

	card := boxStyle.Render(content)

	w, h := termWidth, termHeight
	if w <= 0 {
		w = DefaultWidth
	}
	if h <= 0 {
		h = DefaultHeight
	}

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, card)
}
