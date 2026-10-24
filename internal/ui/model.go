package ui

import (
	"time"

	"attention-span/internal/model"
	"attention-span/internal/store"

	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	ScreenSetup Screen = iota
	ScreenTimer
	ScreenSummary
	ScreenHistory
)

type tickMsg time.Time

// Model is the primary Bubble Tea model managing attention span tracking.
type Model struct {
	screen         Screen
	session        *model.Session
	targetDuration time.Duration
	store          *store.Store
	currentTime    time.Time
	width          int
	height         int
	previousScreen Screen
	cachedHistory  []model.Session
	cachedStats    *store.HistoryStats
}

// NewModel instantiates the UI Model.
func NewModel(st *store.Store, initialTarget time.Duration) Model {
	if initialTarget <= 0 {
		initialTarget = 2 * time.Hour
	}
	return Model{
		screen:         ScreenSetup,
		targetDuration: initialTarget,
		store:          st,
		currentTime:    time.Now(),
		width:          DefaultWidth,
		height:         DefaultHeight,
	}
}

func doTick() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.currentTime = time.Time(msg)
		return m, doTick()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.session != nil && m.session.State != model.StateCompleted && m.session.State != model.StateSetup {
				m.session.Finish(m.currentTime)
				if m.store != nil {
					_ = m.store.SaveSession(m.session)
				}
			}
			return m, tea.Quit
		}

		switch m.screen {
		case ScreenSetup:
			return m.handleSetupKeys(msg)
		case ScreenTimer:
			return m.handleTimerKeys(msg)
		case ScreenSummary:
			return m.handleSummaryKeys(msg)
		case ScreenHistory:
			return m.handleHistoryKeys(msg)
		}
	}

	return m, nil
}

// handleSetupKeys processes keyboard interactions on the duration setup screen.
func (m Model) handleSetupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "up":
		m.targetDuration += 1 * time.Hour
		if m.targetDuration > 12*time.Hour {
			m.targetDuration = 12 * time.Hour
		}

	case "down":
		if m.targetDuration > 1*time.Hour {
			m.targetDuration -= 1 * time.Hour
		} else if m.targetDuration > 15*time.Minute {
			m.targetDuration -= 15 * time.Minute
		}

	case "right":
		m.targetDuration += 15 * time.Minute
		if m.targetDuration > 12*time.Hour {
			m.targetDuration = 12 * time.Hour
		}

	case "left":
		if m.targetDuration > 15*time.Minute {
			m.targetDuration -= 15 * time.Minute
		} else if m.targetDuration > 5*time.Minute {
			m.targetDuration = 5 * time.Minute
		}

	case "h":
		m.previousScreen = ScreenSetup
		m.screen = ScreenHistory
		m.loadHistory()

	case "enter", " ":
		m.session = model.NewSession(m.targetDuration)
		m.session.Start(m.currentTime)
		m.screen = ScreenTimer
	}

	return m, nil
}

// handleTimerKeys processes keyboard interactions during an active focus session.
func (m Model) handleTimerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		if m.session != nil {
			m.session.Finish(m.currentTime)
			if m.store != nil {
				_ = m.store.SaveSession(m.session)
			}
		}
		return m, tea.Quit

	case "e":
		// End session early
		if m.session != nil {
			m.session.Finish(m.currentTime)
			if m.store != nil {
				_ = m.store.SaveSession(m.session)
			}
			m.screen = ScreenSummary
		}

	case " ", "d":
		// Toggle attention / distraction state
		if m.session == nil {
			return m, nil
		}
		if m.session.State == model.StateFocused {
			m.session.HoldAttention(m.currentTime)
		} else if m.session.State == model.StateDistracted {
			m.session.ResumeAttention(m.currentTime)
		}
	}

	return m, nil
}

// handleSummaryKeys processes keyboard interactions on the post-session summary screen.
func (m Model) handleSummaryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case " ":
		m.screen = ScreenSetup
		m.session = nil

	case "h":
		m.previousScreen = ScreenSummary
		m.screen = ScreenHistory
		m.loadHistory()
	}

	return m, nil
}

// handleHistoryKeys processes keyboard interactions on the session history screen.
func (m Model) handleHistoryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "esc":
		m.screen = m.previousScreen
	}

	return m, nil
}

// loadHistory fetches recorded sessions and aggregate lifetime stats from store.
func (m *Model) loadHistory() {
	if m.store != nil {
		m.cachedHistory, _ = m.store.LoadSessions()
		m.cachedStats, _ = m.store.GetStats()
	}
}

// View delegates rendering to the active screen view.
func (m Model) View() string {
	switch m.screen {
	case ScreenSetup:
		return RenderSetupView(m.targetDuration, m.width, m.height)
	case ScreenTimer:
		if m.session != nil {
			return RenderTimerView(m.session, m.currentTime, m.width, m.height)
		}
		return RenderSetupView(m.targetDuration, m.width, m.height)
	case ScreenSummary:
		if m.session != nil {
			return RenderSummaryView(m.session, m.width, m.height)
		}
		return RenderSetupView(m.targetDuration, m.width, m.height)
	case ScreenHistory:
		return RenderHistoryView(m.cachedStats, m.cachedHistory, m.width, m.height)
	default:
		return ""
	}
}
