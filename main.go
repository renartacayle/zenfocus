package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renartacayle/zenfocus/internal/timer"
	"github.com/renartacayle/zenfocus/internal/ui"
)

var (
	version = "1.0.0"
	commit  = "none"
	date    = "unknown"
)

func main() {
	focusMin := flag.Int("focus", 25, "Focus duration in minutes")
	shortBreakMin := flag.Int("break", 5, "Short break duration in minutes")
	longBreakMin := flag.Int("long", 15, "Long break duration in minutes")
	interval := flag.Int("interval", 4, "Number of focus sessions before a long break")
	notify := flag.Bool("notify", true, "Enable desktop notifications (notify-send)")
	showVersion := flag.Bool("v", false, "Print version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "⚡ zenfocus — Minimalist Aesthetic Pomodoro TUI\n\n")
		fmt.Fprintf(os.Stderr, "Usage: zenfocus [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("zenfocus v%s (%s, built at %s)\n", version, commit, date)
		os.Exit(0)
	}

	cfg := timer.Config{
		FocusDuration:      time.Duration(*focusMin) * time.Minute,
		ShortBreakDuration: time.Duration(*shortBreakMin) * time.Minute,
		LongBreakDuration:  time.Duration(*longBreakMin) * time.Minute,
		LongBreakInterval:  *interval,
	}

	state := timer.NewState(cfg)
	state.NotificationsOn = *notify

	model := ui.NewModel(state)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running zenfocus: %v\n", err)
		os.Exit(1)
	}
}
