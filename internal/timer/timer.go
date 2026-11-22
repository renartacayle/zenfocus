package timer

import (
	"fmt"
	"os/exec"
	"time"
)

// SessionType defines the current mode of the timer
type SessionType int

const (
	SessionFocus SessionType = iota
	SessionShortBreak
	SessionLongBreak
)

func (s SessionType) String() string {
	switch s {
	case SessionFocus:
		return "FOCUS SESSION"
	case SessionShortBreak:
		return "SHORT BREAK"
	case SessionLongBreak:
		return "LONG BREAK"
	default:
		return "UNKNOWN"
	}
}

// Config stores durations for each session type
type Config struct {
	FocusDuration      time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
	LongBreakInterval  int // number of focus sessions before a long break
}

// DefaultConfig returns recommended standard Pomodoro timing
func DefaultConfig() Config {
	return Config{
		FocusDuration:      25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		LongBreakInterval:  4,
	}
}

// State represents the runtime state of the timer
type State struct {
	Config           Config
	CurrentSession   SessionType
	RemainingTime    time.Duration
	TotalDuration    time.Duration
	IsRunning        bool
	CompletedCycles  int
	TotalFocusTime   time.Duration
	NotificationsOn  bool
}

// NewState initializes a new timer state
func NewState(cfg Config) *State {
	return &State{
		Config:           cfg,
		CurrentSession:   SessionFocus,
		RemainingTime:    cfg.FocusDuration,
		TotalDuration:    cfg.FocusDuration,
		IsRunning:        false,
		CompletedCycles:  0,
		TotalFocusTime:   0,
		NotificationsOn:  true,
	}
}

// ToggleRunning pauses or resumes the timer
func (s *State) ToggleRunning() {
	s.IsRunning = !s.IsRunning
}

// Reset resets the current session to its initial duration
func (s *State) Reset() {
	s.IsRunning = false
	s.RemainingTime = s.TotalDuration
}

// SwitchSession manually changes the session type
func (s *State) SwitchSession(target SessionType) {
	s.IsRunning = false
	s.CurrentSession = target

	switch target {
	case SessionFocus:
		s.TotalDuration = s.Config.FocusDuration
	case SessionShortBreak:
		s.TotalDuration = s.Config.ShortBreakDuration
	case SessionLongBreak:
		s.TotalDuration = s.Config.LongBreakDuration
	}
	s.RemainingTime = s.TotalDuration
}

// AdvanceSession finishes the current session and rotates to the next state
func (s *State) AdvanceSession() {
	s.IsRunning = false

	if s.CurrentSession == SessionFocus {
		s.CompletedCycles++
		s.TotalFocusTime += s.Config.FocusDuration
		s.notify("Pomodoro Complete!", "Great job! Time to take a restful break.")

		if s.CompletedCycles%s.Config.LongBreakInterval == 0 {
			s.SwitchSession(SessionLongBreak)
		} else {
			s.SwitchSession(SessionShortBreak)
		}
	} else {
		s.notify("Break Over!", "Ready to dive back into deep focus?")
		s.SwitchSession(SessionFocus)
	}
}

// Tick decrements the timer by one second
func (s *State) Tick() bool {
	if !s.IsRunning {
		return false
	}

	if s.RemainingTime > time.Second {
		s.RemainingTime -= time.Second
		return false
	}

	// Session finished
	s.RemainingTime = 0
	s.AdvanceSession()
	return true
}

// Progress calculates the percentage completed (0.0 to 1.0)
func (s *State) Progress() float64 {
	if s.TotalDuration == 0 {
		return 0
	}
	elapsed := s.TotalDuration - s.RemainingTime
	return float64(elapsed) / float64(s.TotalDuration)
}

// FormattedRemaining returns "MM:SS"
func (s *State) FormattedRemaining() string {
	totalSec := int(s.RemainingTime.Seconds())
	mins := totalSec / 60
	secs := totalSec % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

// notify sends a desktop notification if notify-send is available
func (s *State) notify(title, message string) {
	if !s.NotificationsOn {
		return
	}
	_ = exec.Command("notify-send", "-a", "zenfocus", title, message, "-i", "alarm-clock").Start()
}
