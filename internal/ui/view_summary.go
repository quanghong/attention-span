package ui

import (
	"fmt"
	"strings"

	"attention-span/internal/model"

	"github.com/charmbracelet/lipgloss"
)

// RenderSummaryView displays the final post-session breakdown.
func RenderSummaryView(sess *model.Session, width, height int) string {
	boxWidth := BoxWidth(width)

	title := TitleStyle.Render("📊 SESSION COMPLETED")

	focusRatio := sess.FocusRatio * 100
	var gradeBadge string
	if focusRatio >= 85 {
		gradeBadge = lipgloss.NewStyle().Bold(true).Foreground(ColorFocus).Render("🌟 EXCELLENT DEEP WORK FLOW")
	} else if focusRatio >= 70 {
		gradeBadge = lipgloss.NewStyle().Bold(true).Foreground(ColorElapsed).Render("⚡ SOLID CONCENTRATION")
	} else if focusRatio >= 50 {
		gradeBadge = lipgloss.NewStyle().Bold(true).Foreground(ColorHighlight).Render("⚠️ MODERATE DISTRACTION")
	} else {
		gradeBadge = lipgloss.NewStyle().Bold(true).Foreground(ColorDistract).Render("🚨 HEAVY DISTRACTION INTERRUPTIONS")
	}

	actualElapsed := sess.EndTime.Sub(sess.StartTime)

	tableStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	statsList := []string{
		fmt.Sprintf("• Target Planned:    %s", FormatDuration(sess.TargetDuration)),
		fmt.Sprintf("• Actual Session:    %s", FormatDuration(actualElapsed)),
		fmt.Sprintf("• Pure Focus Time:   %s", FormatDuration(sess.TotalFocusTime)),
		fmt.Sprintf("• Distracted Time:   %s", FormatDuration(sess.TotalDistractTime)),
		fmt.Sprintf("• Focus Score:       %.1f%%", focusRatio),
		fmt.Sprintf("• Distraction Drops: %d recorded", len(sess.Drops)),
	}
	statsContent := tableStyle.Render(strings.Join(statsList, "\n"))

	// List drops if any (limit to last 3 so it fits vertically)
	var dropsSection string
	if len(sess.Drops) > 0 {
		dropsHeader := lipgloss.NewStyle().Bold(true).Foreground(ColorAccent).Render("\nDistraction Drop Log:")
		var dropLines []string
		startIdx := 0
		if len(sess.Drops) > 3 {
			startIdx = len(sess.Drops) - 3
		}
		for i := startIdx; i < len(sess.Drops); i++ {
			drop := sess.Drops[i]
			startStr := drop.StartTime.Format("03:04:05 PM")
			durStr := FormatDuration(drop.Duration)
			dropLines = append(dropLines, fmt.Sprintf("  #%d: at %s (lost %s)", i+1, startStr, durStr))
		}
		dropsSection = dropsHeader + "\n" + strings.Join(dropLines, "\n")
	}

	controls := MetricLabel.Render(
		"[Space] Start New Session   │   [h] History   │   [q] Quit",
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		gradeBadge,
		"",
		statsContent,
		dropsSection,
		"",
		controls,
	)

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
