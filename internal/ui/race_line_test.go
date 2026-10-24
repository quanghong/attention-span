package ui

import (
	"strings"
	"testing"
	"time"

	"attention-span/internal/model"
)

// TestFormatDuration verifies that the FormatDuration function correctly formats
// various durations into human-readable strings.
func TestFormatDuration(t *testing.T) {
	cases := []struct {
		d        time.Duration
		expected string
	}{
		{45 * time.Second, "00m 45s"},
		{15 * time.Minute, "15m 00s"},
		{1*time.Hour + 30*time.Minute, "01h 30m 00s"},
		{2*time.Hour + 5*time.Minute + 12*time.Second, "02h 05m 12s"},
	}

	for _, c := range cases {
		got := FormatDuration(c.d)
		if got != c.expected {
			t.Errorf("FormatDuration(%v) = %s, expected %s", c.d, got, c.expected)
		}
	}
}

// TestRenderRaceLines verifies that the race lines for the session are rendered
// correctly with appropriate headers.
func TestRenderRaceLines(t *testing.T) {
	sess := model.NewSession(1 * time.Hour)
	now := time.Now()
	sess.Start(now)

	// Line 1 render
	line1 := RenderDurationLine(sess, now.Add(30*time.Minute), 40)
	if !strings.Contains(line1, "Target Session Line") {
		t.Errorf("line1 missing header: %s", line1)
	}

	// Add a distraction drop
	sess.HoldAttention(now.Add(10 * time.Minute))
	sess.ResumeAttention(now.Add(20 * time.Minute))

	// Line 2 render
	line2 := RenderAttentionRaceLine(sess, now.Add(30*time.Minute), 40)
	if !strings.Contains(line2, "Attention Race") {
		t.Errorf("line2 missing header: %s", line2)
	}
}
