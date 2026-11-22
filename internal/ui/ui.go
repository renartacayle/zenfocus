package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/renartacayle/zenfocus/internal/timer"
)

// Catppuccin Mocha Color Palette
var (
	colorMauve    = lipgloss.Color("#cba6f7")
	colorLavender = lipgloss.Color("#b4befe")
	colorBlue     = lipgloss.Color("#89b4fa")
	colorGreen    = lipgloss.Color("#a6e3a1")
	colorYellow   = lipgloss.Color("#f9e2af")
	colorRed      = lipgloss.Color("#f38ba8")
	colorText     = lipgloss.Color("#cdd6f4")
	colorSubtext  = lipgloss.Color("#a6adc8")
	colorOverlay  = lipgloss.Color("#45475a")
	colorSurface  = lipgloss.Color("#313244")
	colorBase     = lipgloss.Color("#1e1e2e")
)

// Styles
var (
	windowStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMauve).
			Padding(1, 3).
			Align(lipgloss.Center)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorMauve).
			MarginBottom(1)

	badgeFocusStyle = lipgloss.NewStyle().
			Bold(true).
			Background(colorMauve).
			Foreground(colorBase).
			Padding(0, 1).
			MarginBottom(1)

	badgeBreakStyle = lipgloss.NewStyle().
			Bold(true).
			Background(colorGreen).
			Foreground(colorBase).
			Padding(0, 1).
			MarginBottom(1)

	timerDigitsStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorText).
				MarginBottom(1)

	statusRunning = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGreen)

	statusPaused = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorSubtext).
			MarginTop(1)

	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBlue)
)

// ASCII representation for large clock digits (height 3)
var bigDigits = map[rune][]string{
	'0': {"█▀█", "█ █", "▀▀▀"},
	'1': {" █ ", " █ ", " ▀ "},
	'2': {"▀▀█", "█▀▀", "▀▀▀"},
	'3': {"▀▀█", " ▀█", "▀▀▀"},
	'4': {"█ █", "▀▀█", "  ▀"},
	'5': {"█▀▀", "▀▀█", "▀▀▀"},
	'6': {"█▀▀", "█▀█", "▀▀▀"},
	'7': {"▀▀█", "  █", "  ▀"},
	'8': {"█▀█", "█▀█", "▀▀▀"},
	'9': {"█▀█", "▀▀█", "▀▀▀"},
	':': {" ", "■", "■"},
}

// TickMsg is delivered every second
type TickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Model defines the Bubbletea application model
type Model struct {
	state       *timer.State
	progressBar progress.Model
	width       int
	height      int
	quitting    bool
}

// NewModel instantiates the UI Model
func NewModel(state *timer.State) Model {
	p := progress.New(
		progress.WithSolidFill("#cba6f7"),
		progress.WithoutPercentage(),
		progress.WithWidth(32),
	)
	return Model{
		state:       state,
		progressBar: p,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case " ":
			m.state.ToggleRunning()
		case "r":
			m.state.Reset()
		case "s":
			m.state.AdvanceSession()
		case "f":
			m.state.SwitchSession(timer.SessionFocus)
		case "b":
			m.state.SwitchSession(timer.SessionShortBreak)
		case "l":
			m.state.SwitchSession(timer.SessionLongBreak)
		}

	case TickMsg:
		m.state.Tick()
		return m, tickCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// renderBigTime formats "MM:SS" into high-impact 3-line ASCII numbers
func renderBigTime(timeStr string) string {
	lines := []string{"", "", ""}
	for _, ch := range timeStr {
		digitArt, exists := bigDigits[ch]
		if !exists {
			digitArt = []string{" ", " ", " "}
		}
		for i := 0; i < 3; i++ {
			lines[i] += digitArt[i] + " "
		}
	}
	return strings.Join(lines, "\n")
}

func (m Model) View() string {
	if m.quitting {
		return lipgloss.NewStyle().Foreground(colorMauve).Render("\n  ✨ Thank you for using zenfocus. Stay sharp & productive!\n\n")
	}

	// 1. Title Banner
	title := titleStyle.Render("⚡ Z E N F O C U S")

	// 2. Mode Badge
	var badge string
	if m.state.CurrentSession == timer.SessionFocus {
		badge = badgeFocusStyle.Render("  🍅 FOCUS SESSION  ")
	} else if m.state.CurrentSession == timer.SessionShortBreak {
		badge = badgeBreakStyle.Render("  ☕ SHORT BREAK  ")
	} else {
		badge = badgeBreakStyle.Render("  🌴 LONG BREAK  ")
	}

	// 3. Status Indicator
	var status string
	if m.state.IsRunning {
		status = statusRunning.Render("▶ RUNNING")
	} else {
		status = statusPaused.Render("⏸ PAUSED")
	}

	// 4. Large ASCII Timer
	bigClock := timerDigitsStyle.Render(renderBigTime(m.state.FormattedRemaining()))

	// 5. Progress Bar
	progressBar := m.progressBar.ViewAs(m.state.Progress())

	// 6. Session & Cycles Counter
	cycleIndex := (m.state.CompletedCycles % m.state.Config.LongBreakInterval)
	var tomatoes string
	for i := 0; i < m.state.Config.LongBreakInterval; i++ {
		if i < cycleIndex {
			tomatoes += "🍅 "
		} else {
			tomatoes += "○  "
		}
	}
	cycleInfo := lipgloss.NewStyle().Foreground(colorSubtext).Render(
		fmt.Sprintf("Cycles: %s (%d completed)", tomatoes, m.state.CompletedCycles),
	)

	// 7. Navigation & Controls Help
	controls := helpStyle.Render(
		fmt.Sprintf("%s toggle  •  %s skip  •  %s reset  •  %s focus  •  %s break  •  %s quit",
			keyStyle.Render("[Space]"),
			keyStyle.Render("[s]"),
			keyStyle.Render("[r]"),
			keyStyle.Render("[f]"),
			keyStyle.Render("[b]"),
			keyStyle.Render("[q]"),
		),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		badge,
		status,
		"",
		bigClock,
		progressBar,
		"",
		cycleInfo,
		controls,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		windowStyle.Render(content),
	)
}
