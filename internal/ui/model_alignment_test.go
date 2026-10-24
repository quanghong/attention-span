package ui

import (
	"strings"
	"testing"
	"time"

	"attention-span/internal/model"
)

// TestVisualAlignment verifies that the rendered setup and timer views align correctly
// within the specified dimensions and do not contain unexpected ASCII block characters.
func TestVisualAlignment(t *testing.T) {
	setupView := RenderSetupView(2*time.Hour, 95, 20)
	t.Logf("\nRendered Setup View (95x20):\n%s\n", setupView)

	sess := model.NewSession(2 * time.Hour)
	now := time.Now()
	sess.Start(now)
	timerView := RenderTimerView(sess, now, 95, 20)
	t.Logf("\nRendered Timer View (95x20):\n%s\n", timerView)

	// Ensure no BigClock ASCII block art is present
	if strings.Contains(timerView, "|_|") || strings.Contains(setupView, "|_|") {
		t.Error("found unexpected big block characters")
	}

	// Verify that height fits nicely into 20 rows
	setupLines := strings.Split(setupView, "\n")
	if len(setupLines) > 20 {
		t.Errorf("setupView exceeds 20 lines: got %d", len(setupLines))
	}

	timerLines := strings.Split(timerView, "\n")
	if len(timerLines) > 20 {
		t.Errorf("timerView exceeds 20 lines: got %d", len(timerLines))
	}
}
