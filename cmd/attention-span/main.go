package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"attention-span/internal/store"
	"attention-span/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	version = "1.0.0"
)

func main() {
	durationFlag := flag.String("duration", "2h", "Target duration for focus session (e.g. 2h, 45m, 1h30m)")
	flag.StringVar(durationFlag, "d", "2h", "Target duration (shorthand)")

	historyFlag := flag.Bool("history", false, "Display past session history and exit")
	statsFlag := flag.Bool("stats", false, "Display lifetime focus statistics and exit")
	versionFlag := flag.Bool("version", false, "Print application version")
	flag.BoolVar(versionFlag, "v", false, "Print application version (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "⚡ Attention Span Tracker v%s ⚡\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n  attention-span [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nControls in TUI:\n")
		fmt.Fprintf(os.Stderr, "  [Space] / [d]  Flag distraction (hold attention)\n")
		fmt.Fprintf(os.Stderr, "  [Space] / [c]  Resume concentration (continue)\n")
		fmt.Fprintf(os.Stderr, "  [e]            End session and view summary\n")
		fmt.Fprintf(os.Stderr, "  [q]            Quit\n")
	}

	flag.Parse()

	if *versionFlag {
		fmt.Printf("attention-span v%s\n", version)
		os.Exit(0)
	}

	// Initialize storage
	sessionStore, err := store.NewStore("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not initialize storage (%v). Sessions will not be persisted.\n", err)
	}

	if *historyFlag {
		if sessionStore == nil {
			fmt.Println("No store available.")
			os.Exit(0)
		}
		sessions, err := sessionStore.LoadSessions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading history: %v\n", err)
			os.Exit(1)
		}
		if len(sessions) == 0 {
			fmt.Println("No past sessions found.")
			os.Exit(0)
		}
		fmt.Println("Recent Focus Sessions:")
		for i, s := range sessions {
			fmt.Printf("  #%d [%s] Target: %s | Focused: %s | Ratio: %.1f%% | Drops: %d\n",
				i+1,
				s.StartTime.Format("2006-01-02 15:04"),
				ui.FormatDurationShort(s.TargetDuration),
				ui.FormatDurationShort(s.TotalFocusTime),
				s.FocusRatio*100,
				len(s.Drops),
			)
		}
		os.Exit(0)
	}

	if *statsFlag {
		if sessionStore == nil {
			fmt.Println("No store available.")
			os.Exit(0)
		}
		stats, err := sessionStore.GetStats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading stats: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Lifetime Attention Span Statistics:\n")
		fmt.Printf("  Total Sessions:        %d\n", stats.TotalSessions)
		fmt.Printf("  Total Planned Time:    %s\n", ui.FormatDuration(stats.TotalPlannedTime))
		fmt.Printf("  Total Focus Time:      %s\n", ui.FormatDuration(stats.TotalFocusTime))
		fmt.Printf("  Total Distracted Time: %s\n", ui.FormatDuration(stats.TotalDistractedTime))
		fmt.Printf("  Total Drops:           %d\n", stats.TotalDrops)
		fmt.Printf("  Average Focus Score:   %.1f%%\n", stats.AverageFocusRatio*100)
		os.Exit(0)
	}

	// Parse target duration
	target, err := time.ParseDuration(*durationFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid duration %q, defaulting to 2h\n", *durationFlag)
		target = 2 * time.Hour
	}

	m := ui.NewModel(sessionStore, target)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running attention-span: %v\n", err)
		os.Exit(1)
	}
}
