package model

import (
	"fmt"
	"time"
)

// SessionState defines the active lifecycle state of a tracking session.
type SessionState int

const (
	StateSetup SessionState = iota
	StateFocused
	StateDistracted
	StatePaused
	StateCompleted
)

func (s SessionState) String() string {
	switch s {
	case StateSetup:
		return "Setup"
	case StateFocused:
		return "Focused"
	case StateDistracted:
		return "Distracted"
	case StatePaused:
		return "Paused"
	case StateCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

// DistractionDrop represents a single distraction event where attention dropped.
type DistractionDrop struct {
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Reason    string        `json:"reason,omitempty"`
}

// CalculateDuration computes duration, handling ongoing drops where EndTime is not yet set.
func (d DistractionDrop) CalculateDuration(now time.Time) time.Duration {
	if !d.EndTime.IsZero() {
		return d.EndTime.Sub(d.StartTime)
	}
	if !d.StartTime.IsZero() && now.After(d.StartTime) {
		return now.Sub(d.StartTime)
	}
	return 0
}

// Session records the full timeline, targets, and distraction drops of a work block.
type Session struct {
	ID                string            `json:"id"`
	TargetDuration    time.Duration     `json:"target_duration"`
	StartTime         time.Time         `json:"start_time"`
	EndTime           time.Time         `json:"end_time"`
	Drops             []DistractionDrop `json:"drops"`
	CurrentDrop       *DistractionDrop  `json:"current_drop,omitempty"`
	TotalFocusTime    time.Duration     `json:"total_focus_time"`
	TotalDistractTime time.Duration     `json:"total_distract_time"`
	FocusRatio        float64           `json:"focus_ratio"`
	State             SessionState      `json:"state"`
}

// NewSession creates an initialized session with the given target duration.
func NewSession(target time.Duration) *Session {
	return &Session{
		ID:             fmt.Sprintf("sess_%d", time.Now().UnixNano()),
		TargetDuration: target,
		Drops:          make([]DistractionDrop, 0),
		State:          StateSetup,
	}
}

// Start marks the beginning of the session.
func (s *Session) Start(now time.Time) {
	if s.State != StateSetup {
		return
	}
	s.StartTime = now
	s.State = StateFocused
}

// HoldAttention flags that the user became distracted ("hold attention").
func (s *Session) HoldAttention(now time.Time) {
	if s.State != StateFocused {
		return
	}
	s.State = StateDistracted
	s.CurrentDrop = &DistractionDrop{
		StartTime: now,
	}
}

// ResumeAttention flags that concentration has resumed ("continue").
func (s *Session) ResumeAttention(now time.Time) {
	if s.State != StateDistracted {
		return
	}
	if s.CurrentDrop != nil {
		s.CurrentDrop.EndTime = now
		s.CurrentDrop.Duration = now.Sub(s.CurrentDrop.StartTime)
		s.Drops = append(s.Drops, *s.CurrentDrop)
		s.CurrentDrop = nil
	}
	s.State = StateFocused
}

// Finish concludes the session and computes final metrics.
func (s *Session) Finish(now time.Time) {
	if s.State == StateCompleted {
		return
	}
	if s.State == StateDistracted && s.CurrentDrop != nil {
		s.CurrentDrop.EndTime = now
		s.CurrentDrop.Duration = now.Sub(s.CurrentDrop.StartTime)
		s.Drops = append(s.Drops, *s.CurrentDrop)
		s.CurrentDrop = nil
	}
	s.EndTime = now
	s.State = StateCompleted
	s.TotalDistractTime = s.TotalDistractionTime(now)
	s.TotalFocusTime = s.CalculateFocusTime(now)
	s.FocusRatio = s.CalculateFocusRatio(now)
}

// Elapsed returns total session wall-clock time from start until now (or end).
func (s *Session) Elapsed(now time.Time) time.Duration {
	if s.StartTime.IsZero() {
		return 0
	}
	if !s.EndTime.IsZero() {
		return s.EndTime.Sub(s.StartTime)
	}
	if now.Before(s.StartTime) {
		return 0
	}
	return now.Sub(s.StartTime)
}

// Remaining returns how much target time is left.
func (s *Session) Remaining(now time.Time) time.Duration {
	elapsed := s.Elapsed(now)
	if elapsed >= s.TargetDuration {
		return 0
	}
	return s.TargetDuration - elapsed
}

// TotalDistractionTime returns the cumulative duration of all past and current drops.
func (s *Session) TotalDistractionTime(now time.Time) time.Duration {
	var total time.Duration
	for _, drop := range s.Drops {
		if !drop.EndTime.IsZero() {
			total += drop.Duration
		} else {
			total += drop.CalculateDuration(now)
		}
	}
	if s.CurrentDrop != nil {
		total += s.CurrentDrop.CalculateDuration(now)
	}
	return total
}

// CalculateFocusTime returns elapsed wall-clock time minus total distraction time.
func (s *Session) CalculateFocusTime(now time.Time) time.Duration {
	elapsed := s.Elapsed(now)
	distracted := s.TotalDistractionTime(now)
	if distracted >= elapsed {
		return 0
	}
	return elapsed - distracted
}

// CalculateFocusRatio returns the percentage (0.0 to 1.0) of elapsed time spent focused.
func (s *Session) CalculateFocusRatio(now time.Time) float64 {
	elapsed := s.Elapsed(now)
	if elapsed <= 0 {
		return 1.0
	}
	focus := s.CalculateFocusTime(now)
	ratio := float64(focus) / float64(elapsed)
	if ratio < 0 {
		return 0
	}
	if ratio > 1 {
		return 1.0
	}
	return ratio
}

// TargetProgress returns the percentage (0.0 to 1.0) of target duration reached.
func (s *Session) TargetProgress(now time.Time) float64 {
	if s.TargetDuration <= 0 {
		return 1.0
	}
	elapsed := s.Elapsed(now)
	progress := float64(elapsed) / float64(s.TargetDuration)
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

// FocusProgress returns the percentage (0.0 to 1.0) of target duration spent focused.
func (s *Session) FocusProgress(now time.Time) float64 {
	if s.TargetDuration <= 0 {
		return 1.0
	}
	focus := s.CalculateFocusTime(now)
	progress := float64(focus) / float64(s.TargetDuration)
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

// IsTargetReached returns true if elapsed duration has met or exceeded the target.
func (s *Session) IsTargetReached(now time.Time) bool {
	if s.StartTime.IsZero() || s.TargetDuration <= 0 {
		return false
	}
	return s.Elapsed(now) >= s.TargetDuration
}
