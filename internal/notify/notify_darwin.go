//go:build darwin

package notify

import (
	"fmt"
	"os/exec"
)

// Send posts a native notification to the macOS Notification Center.
func Send(title, message string) error {
	script := fmt.Sprintf(`display notification %q with title %q sound name "Ping"`, message, title)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}
