package ui

import (
	"fmt"
	"strings"

	"attention-span/internal/model"
	"attention-span/internal/store"
	"github.com/charmbracelet/lipgloss"
)

// RenderHistoryView displays lifetime statistics and recent focus sessions.
func RenderHistoryView(stats *store.HistoryStats, sessions []model.Session, width, height int) string {
	boxWidth := 68
	if width > 0 && width < 72 {
		boxWidth = width - 4
	}

	title := TitleStyle.Render("📜 FOCUS SESSIONS HISTORY")

	if stats == nil || stats.TotalSessions == 0 {
		emptyMsg := MetricLabel.Render("No past sessions found yet. Complete a session to see history!")
		controls := MetricLabel.Render("[Esc / Backspace / h] Back to Setup   │   [q] Quit")
		content := lipgloss.JoinVertical(lipgloss.Center, title, "", emptyMsg, "", controls)
		card := SquareClockBox.Copy().Width(boxWidth).Render(content)
		if height > 0 {
			return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
		}
		return card
	}

	statsBar := fmt.Sprintf(
		"Sessions: %s  │  Focus Time: %s  │  Avg Ratio: %s  │  Drops: %s",
		MetricValue.Render(fmt.Sprintf("%d", stats.TotalSessions)),
		MetricAccent.Render(FormatDurationShort(stats.TotalFocusTime)),
		MetricAccent.Render(fmt.Sprintf("%.1f%%", stats.AverageFocusRatio*100)),
		MetricValue.Render(fmt.Sprintf("%d", stats.TotalDrops)),
	)

	listHeader := lipgloss.NewStyle().Bold(true).Foreground(ColorAccent).Render("\nRecent Sessions:")
	var sessionLines []string
	maxShow := 8
	if len(sessions) < maxShow {
		maxShow = len(sessions)
	}

	for i := 0; i < maxShow; i++ {
		sess := sessions[i]
		dateStr := sess.StartTime.Format("01/02 15:04")
		ratioStr := fmt.Sprintf("%.0f%%", sess.FocusRatio*100)
		sessionLines = append(sessionLines, fmt.Sprintf(
			"  %s │ Target: %-6s │ Focused: %-6s │ Ratio: %-4s │ Drops: %d",
			dateStr,
			FormatDurationShort(sess.TargetDuration),
			FormatDurationShort(sess.TotalFocusTime),
			ratioStr,
			len(sess.Drops),
		))
	}

	historyList := listHeader + "\n" + strings.Join(sessionLines, "\n")

	controls := MetricLabel.Render(
		"[Esc / Backspace / h] Back to Setup   │   [q] Quit",
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		statsBar,
		historyList,
		"",
		controls,
	)

	card := SquareClockBox.Copy().Width(boxWidth).Render(content)

	if height > 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, card)
	}

	return card
}
