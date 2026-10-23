//go:build !darwin

package notify

import (
	"os/exec"
)

// Send attempts to post a desktop notification using notify-send on Linux or falls back silently.
func Send(title, message string) error {
	if path, err := exec.LookPath("notify-send"); err == nil && path != "" {
		_ = exec.Command("notify-send", title, message).Run()
	}
	return nil
}
