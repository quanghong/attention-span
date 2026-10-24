package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"attention-span/internal/model"
)

// HistoryStats aggregates all-time focus statistics.
type HistoryStats struct {
	TotalSessions       int           `json:"total_sessions"`
	TotalPlannedTime    time.Duration `json:"total_planned_time"`
	TotalFocusTime      time.Duration `json:"total_focus_time"`
	TotalDistractedTime time.Duration `json:"total_distracted_time"`
	TotalDrops          int           `json:"total_drops"`
	AverageFocusRatio   float64       `json:"average_focus_ratio"`
}

// Store handles persisting session history to disk.
type Store struct {
	mu       sync.RWMutex
	filePath string
}

// DefaultPath returns the standard platform config path for attention-span history.
func DefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, hErr := os.UserHomeDir()
		if hErr != nil {
			return "", errors.New("unable to determine user config or home directory")
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "attention-span", "history.json"), nil
}

// NewStore creates a new Store instance. If path is empty, default path is used.
func NewStore(customPath string) (*Store, error) {
	path := customPath
	if path == "" {
		defPath, err := DefaultPath()
		if err != nil {
			return nil, err
		}
		path = defPath
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory %s: %w", dir, err)
	}

	return &Store{filePath: path}, nil
}

// SaveSession appends or updates a session in the JSON file.
func (s *Store) SaveSession(sess *model.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessions, err := s.readSessionsLocked()
	if err != nil {
		sessions = []model.Session{}
	}

	// Update existing if present, otherwise append
	updated := false
	for i, existing := range sessions {
		if existing.ID == sess.ID {
			sessions[i] = *sess
			updated = true
			break
		}
	}
	if !updated {
		sessions = append(sessions, *sess)
	}

	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode sessions to json: %w", err)
	}

	tmpFile := s.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write tmp file: %w", err)
	}

	if err := os.Rename(tmpFile, s.filePath); err != nil {
		return fmt.Errorf("failed to commit file: %w", err)
	}

	return nil
}

// LoadSessions retrieves all stored sessions, newest first.
func (s *Store) LoadSessions() ([]model.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions, err := s.readSessionsLocked()
	if err != nil {
		return nil, err
	}

	// Reverse to show newest first
	for i, j := 0, len(sessions)-1; i < j; i, j = i+1, j-1 {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	}

	return sessions, nil
}

// GetStats calculates lifetime aggregate metrics across all sessions.
func (s *Store) GetStats() (*HistoryStats, error) {
	sessions, err := s.LoadSessions()
	if err != nil {
		return nil, err
	}

	stats := &HistoryStats{
		TotalSessions: len(sessions),
	}

	if len(sessions) == 0 {
		return stats, nil
	}

	var sumRatio float64
	for _, sess := range sessions {
		stats.TotalPlannedTime += sess.TargetDuration
		stats.TotalFocusTime += sess.TotalFocusTime
		stats.TotalDistractedTime += sess.TotalDistractTime
		stats.TotalDrops += len(sess.Drops)
		sumRatio += sess.FocusRatio
	}

	stats.AverageFocusRatio = sumRatio / float64(len(sessions))
	return stats, nil
}

// readSessionsLocked reads the sessions from the JSON file without acquiring any locks.
func (s *Store) readSessionsLocked() ([]model.Session, error) {
	if _, err := os.Stat(s.filePath); errors.Is(err, os.ErrNotExist) {
		return []model.Session{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage file: %w", err)
	}

	if len(data) == 0 {
		return []model.Session{}, nil
	}

	var sessions []model.Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("corrupt history json: %w", err)
	}

	return sessions, nil
}
