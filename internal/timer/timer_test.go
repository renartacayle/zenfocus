package timer

import (
	"testing"
	"time"
)

func TestNewState(t *testing.T) {
	cfg := DefaultConfig()
	state := NewState(cfg)

	if state.CurrentSession != SessionFocus {
		t.Errorf("Expected initial session to be SessionFocus, got %v", state.CurrentSession)
	}

	if state.RemainingTime != 25*time.Minute {
		t.Errorf("Expected remaining time 25m, got %v", state.RemainingTime)
	}

	if state.IsRunning {
		t.Errorf("Expected timer to start paused")
	}
}

func TestToggleRunning(t *testing.T) {
	state := NewState(DefaultConfig())

	state.ToggleRunning()
	if !state.IsRunning {
		t.Errorf("Expected timer to be running after toggle")
	}

	state.ToggleRunning()
	if state.IsRunning {
		t.Errorf("Expected timer to be paused after second toggle")
	}
}

func TestTickAndAdvance(t *testing.T) {
	cfg := Config{
		FocusDuration:      2 * time.Second,
		ShortBreakDuration: 1 * time.Second,
		LongBreakDuration:  3 * time.Second,
		LongBreakInterval:  2,
	}
	state := NewState(cfg)
	state.NotificationsOn = false // disable notify in test
	state.IsRunning = true

	// Tick 1
	state.Tick()
	if state.RemainingTime != 1*time.Second {
		t.Errorf("Expected remaining time 1s, got %v", state.RemainingTime)
	}

	// Tick 2: Should advance to Short Break
	completed := state.Tick()
	if !completed {
		t.Errorf("Expected session to complete")
	}

	if state.CurrentSession != SessionShortBreak {
		t.Errorf("Expected session to be ShortBreak, got %v", state.CurrentSession)
	}

	if state.CompletedCycles != 1 {
		t.Errorf("Expected 1 completed cycle, got %d", state.CompletedCycles)
	}
}

func TestProgressCalculation(t *testing.T) {
	cfg := Config{
		FocusDuration: 10 * time.Second,
	}
	state := NewState(cfg)
	state.TotalDuration = 10 * time.Second
	state.RemainingTime = 5 * time.Second

	p := state.Progress()
	if p != 0.5 {
		t.Errorf("Expected progress 0.5, got %f", p)
	}
}
