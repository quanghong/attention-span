package model

import (
	"testing"
	"time"
)

// TestSessionLifecycle verifies the full lifecycle of a focus session,
// including state transitions, focus and distraction tracking, and progress calculations.
func TestSessionLifecycle(t *testing.T) {
	target := 2 * time.Hour
	sess := NewSession(target)

	if sess.State != StateSetup {
		t.Fatalf("expected StateSetup, got %v", sess.State)
	}

	t0 := time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)
	sess.Start(t0)

	if sess.State != StateFocused {
		t.Fatalf("expected StateFocused, got %v", sess.State)
	}

	// 30 minutes in, user is focused
	t1 := t0.Add(30 * time.Minute)
	if elapsed := sess.Elapsed(t1); elapsed != 30*time.Minute {
		t.Fatalf("expected elapsed 30m, got %v", elapsed)
	}
	if focus := sess.CalculateFocusTime(t1); focus != 30*time.Minute {
		t.Fatalf("expected focus 30m, got %v", focus)
	}
	if ratio := sess.CalculateFocusRatio(t1); ratio != 1.0 {
		t.Fatalf("expected ratio 1.0, got %f", ratio)
	}

	// Hold attention (distracted) at t1 (9:30)
	sess.HoldAttention(t1)
	if sess.State != StateDistracted {
		t.Fatalf("expected StateDistracted, got %v", sess.State)
	}

	// 10 minutes pass in distraction (9:40)
	t2 := t1.Add(10 * time.Minute)
	if dist := sess.TotalDistractionTime(t2); dist != 10*time.Minute {
		t.Fatalf("expected distracted 10m, got %v", dist)
	}
	if focus := sess.CalculateFocusTime(t2); focus != 30*time.Minute {
		t.Fatalf("expected focus 30m, got %v", focus)
	}

	// Resume attention at t2 (9:40)
	sess.ResumeAttention(t2)
	if sess.State != StateFocused {
		t.Fatalf("expected StateFocused after resume, got %v", sess.State)
	}
	if len(sess.Drops) != 1 {
		t.Fatalf("expected 1 completed drop, got %d", len(sess.Drops))
	}
	if sess.Drops[0].Duration != 10*time.Minute {
		t.Fatalf("expected drop duration 10m, got %v", sess.Drops[0].Duration)
	}

	// Work for another 20 minutes (10:00)
	t3 := t2.Add(20 * time.Minute)
	// Total elapsed: 30 + 10 + 20 = 60 minutes
	// Focus: 50 minutes, Distracted: 10 minutes
	if elapsed := sess.Elapsed(t3); elapsed != 60*time.Minute {
		t.Fatalf("expected elapsed 60m, got %v", elapsed)
	}
	if focus := sess.CalculateFocusTime(t3); focus != 50*time.Minute {
		t.Fatalf("expected focus 50m, got %v", focus)
	}
	expectedRatio := 50.0 / 60.0
	if ratio := sess.CalculateFocusRatio(t3); ratio != expectedRatio {
		t.Fatalf("expected ratio %f, got %f", expectedRatio, ratio)
	}

	// 50m focus out of 2h (120m) target = 50 / 120 = ~0.4166
	if progress := sess.FocusProgress(t3); progress < 0.41 || progress > 0.42 {
		t.Fatalf("expected focus progress ~0.416, got %f", progress)
	}
	// 60m elapsed out of 2h (120m) target = 60 / 120 = 0.5
	if progress := sess.TargetProgress(t3); progress != 0.5 {
		t.Fatalf("expected target progress 0.5, got %f", progress)
	}

	// Finish session at t3
	sess.Finish(t3)
	if sess.State != StateCompleted {
		t.Fatalf("expected StateCompleted, got %v", sess.State)
	}
	if sess.TotalFocusTime != 50*time.Minute {
		t.Fatalf("expected final TotalFocusTime 50m, got %v", sess.TotalFocusTime)
	}
	if sess.TotalDistractTime != 10*time.Minute {
		t.Fatalf("expected final TotalDistractTime 10m, got %v", sess.TotalDistractTime)
	}
}

// TestTargetReached verifies that the session correctly identifies when the target duration has been reached.
func TestTargetReached(t *testing.T) {
	target := 1 * time.Hour
	sess := NewSession(target)
	t0 := time.Now()
	sess.Start(t0)

	if sess.IsTargetReached(t0.Add(59 * time.Minute)) {
		t.Fatal("target should not be reached at 59m")
	}
	if !sess.IsTargetReached(t0.Add(60 * time.Minute)) {
		t.Fatal("target should be reached at 60m")
	}
	if !sess.IsTargetReached(t0.Add(75 * time.Minute)) {
		t.Fatal("target should be reached past 60m")
	}
}
