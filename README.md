# ⚡ zenfocus — Minimalist Aesthetic Pomodoro TUI

[![CI](https://github.com/renartacayle/zenfocus/actions/workflows/ci.yml/badge.svg)](https://github.com/renartacayle/zenfocus/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)
![Theme: Catppuccin Mocha](https://img.shields.io/badge/Theme-Catppuccin%20Mocha-CBA6F7)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)

A lightweight, distraction-free, and keyboard-driven Pomodoro timer for terminal lovers. Built with **Go** and Charmbracelet's **Bubbletea** & **Lip Gloss**, styled natively with the **Catppuccin Mocha** color palette.

---

## 🎯 What is this for?

Most focus and timer apps live in heavy Electron wrappers or browser tabs where notifications and distractions are only a click away.

`zenfocus` brings the Pomodoro technique directly into your terminal workspace:

1. **⌨️ Keyboard-Driven & Distraction-Free**: Runs natively in your terminal, Tmux pane, or tiling window manager. Control everything with single-key shortcuts.
2. **🎨 Catppuccin Mocha Aesthetics**: Beautiful large ASCII clock digits, smooth gradient progress bars, and clean state badges.
3. **🔄 Automated Pomodoro State Machine**: Cycles automatically between **Focus** (25 min), **Short Breaks** (5 min), and **Long Breaks** (15 min) after 4 completed cycles.
4. **🔔 Desktop Notifications**: Integrates with system notification daemons (`notify-send`) so you never miss a break transition while working in full-screen.
5. **⚡ Zero-Dependency Single Binary**: Compiles into a single standalone ~4MB binary that boots in under 5 milliseconds.

---

## 🖥️ Terminal Interface Preview

```text
┌────────────────────────────────────────────────────────────┐
│                    ⚡ Z E N F O C U S                      │
│                   🍅 FOCUS SESSION                         │
│                      ▶ RUNNING                             │
│                                                            │
│              █▀█ █▀▀   ■   █▀█ █▀█                         │
│              █▀▀ ▀▀█       █ █ █ █                         │
│              ▀▀▀ ▀▀▀   ■   ▀▀▀ ▀▀▀                         │
│                                                            │
│        [████████████████████████░░░░░░░░]                  │
│                                                            │
│            Cycles: 🍅 🍅 🍅 ○ (3 completed)                │
│                                                            │
│    [Space] toggle  •  [s] skip  •  [r] reset  •  [q] quit  │
└────────────────────────────────────────────────────────────┘
```

---

## ⌨️ Controls & Keybindings

| Key | Action |
| :--- | :--- |
| `Space` | Toggle timer (**Start / Pause**) |
| `s` | **Skip** current session and advance to the next state |
| `r` | **Reset** current session to original duration |
| `f` | Switch to **Focus Session** mode manually |
| `b` | Switch to **Short Break** mode manually |
| `l` | Switch to **Long Break** mode manually |
| `q` / `Ctrl+C` | **Quit** application cleanly |

---

## 🚀 Installation & Usage

### Option 1: Via Go Install
```bash
go install github.com/renartacayle/zenfocus@latest
```

### Option 2: Build from Source
```bash
git clone https://github.com/renartacayle/zenfocus.git
cd zenfocus
go build -o zenfocus .
./zenfocus
```

### CLI Flags & Custom Durations
You can customize session durations on the fly using command-line flags:

```bash
# 50-minute deep work focus with a 10-minute break
zenfocus -focus 50 -break 10 -long 20

# Disable desktop notifications
zenfocus -notify=false
```

---

## 🧪 Testing

Run the automated unit test suite:
```bash
go test -v ./...
```

---

## 📄 License
This project is open-sourced under the [MIT License](LICENSE).
