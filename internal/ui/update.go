package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all UI interactions.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		b.spinner, cmd = b.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keys.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keys.Exit):
			return b, tea.Quit
		}
	case Message:
		b.messages = append(b.messages[1:], msg)
		return b, cmd
	}

	return b, tea.Batch(cmds...)
}
