package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI.
func (b Bubble) Init() tea.Cmd {
	return b.spinner.Tick
}
