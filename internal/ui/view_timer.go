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
	boxWidth := 66
	if termWidth > 0 && termWidth < 70 {
		boxWidth = termWidth - 4
	}
	trackWidth := boxWidth - 8
	if trackWidth < 20 {
		trackWidth = 20
	}

	// 1. Clock Header
	currentTimeStr := now.Format("03:04:05 PM")
	targetEndTime := sess.StartTime.Add(sess.TargetDuration)
	endTimeStr := targetEndTime.Format("03:04 PM")

	clockHeader := lipgloss.JoinHorizontal(
		lipgloss.Center,
		BigClockStyle.Render("🕒 "+currentTimeStr),
		MetricLabel.Render("  │  Target End: "+endTimeStr),
	)

	// 2. Countdown Remaining
	remaining := sess.Remaining(now)
	countdownStr := RemainingStyle.Render("REMAINING: " + FormatDuration(remaining))

	// 3. Status Badge & Context Message
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
		statusMsg = lipgloss.NewStyle().Foreground(ColorDistract).Render("Press [SPACE] or [c] when you refocus" + currentDropDuration)
	} else if sess.State == model.StateCompleted {
		badge = CompletedBadge.Render("★ TARGET TIME REACHED")
		statusMsg = MetricAccent.Render("Session duration reached! Press [e] to view summary.")
	}

	// 4. Line 1: Duration from Start Time to End Time
	line1 := RenderDurationLine(sess, now, trackWidth)

	// 5. Line 2: Attention Race with Distraction Drops
	line2 := RenderAttentionRaceLine(sess, now, trackWidth)

	// 6. Metrics Row
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

	// 7. Controls Footer
	controls := MetricLabel.Render(
		"[Space] Toggle Attention   │   [e] End Early   │   [q] Quit",
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		clockHeader,
		countdownStr,
		"",
		badge,
		statusMsg,
		"",
		line1,
		"",
		line2,
		"",
		metricsBar,
		"",
		controls,
		"",
	)

	// Style square box according to focus state
	var boxStyle lipgloss.Style
	if sess.State == model.StateDistracted {
		boxStyle = DistractedClockBox.Copy().Width(boxWidth)
	} else {
		boxStyle = SquareClockBox.Copy().Width(boxWidth)
	}

	card := boxStyle.Render(content)

	if termHeight > 0 {
		return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, card)
	}

	return card
}
