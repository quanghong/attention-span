package store

import (
	"path/filepath"
	"testing"
	"time"

	"attention-span/internal/model"
)

// TestStoreSaveAndLoad verifies that sessions can be saved and loaded correctly,
// and that statistics are computed accurately.
func TestStoreSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "history.json")

	st, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Empty load
	sessions, err := st.LoadSessions()
	if err != nil {
		t.Fatalf("failed to load empty sessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("expected 0 sessions, got %d", len(sessions))
	}

	// Create and save session
	s1 := model.NewSession(1 * time.Hour)
	now := time.Now()
	s1.Start(now)
	s1.HoldAttention(now.Add(15 * time.Minute))
	s1.ResumeAttention(now.Add(20 * time.Minute))
	s1.Finish(now.Add(1 * time.Hour))

	if err := st.SaveSession(s1); err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	loaded, err := st.LoadSessions()
	if err != nil {
		t.Fatalf("failed to load sessions: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 session, got %d", len(loaded))
	}
	if loaded[0].ID != s1.ID {
		t.Fatalf("expected ID %s, got %s", s1.ID, loaded[0].ID)
	}
	if len(loaded[0].Drops) != 1 {
		t.Fatalf("expected 1 drop, got %d", len(loaded[0].Drops))
	}

	stats, err := st.GetStats()
	if err != nil {
		t.Fatalf("failed to compute stats: %v", err)
	}
	if stats.TotalSessions != 1 {
		t.Fatalf("expected 1 total session, got %d", stats.TotalSessions)
	}
	if stats.TotalDrops != 1 {
		t.Fatalf("expected 1 total drop, got %d", stats.TotalDrops)
	}
}
