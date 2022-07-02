package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the UI.
func (b Bubble) View() string {
	s := b.spinner.View() + " Receiving events...\n\n"

	for _, msg := range b.messages {
		if msg == nil {
			s += Nil.String() + "\n"
			continue
		}

		s += msg.String() + "\n"
	}

	return lipgloss.NewStyle().Margin(1, 2, 0, 2).Render(s)
}
