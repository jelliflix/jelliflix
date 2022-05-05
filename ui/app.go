package ui

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	err  error
	quit bool

	spinner   spinner.Model
	paginator paginator.Model

	messages []Message
}

func NewApp() (app App) {
	app.spinner = spinner.New()
	app.spinner.Spinner = spinner.Dot
	app.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	app.messages = make([]Message, 10)

	return
}

func (m App) Init() tea.Cmd {
	return spinner.Tick
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quit = true
		return m, tea.Quit
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case Message:
		m.messages = append(m.messages[1:], msg)
		return m, cmd
	}

	return m, cmd
}

func (m App) View() string {
	s := m.spinner.View() + Bold(" Receiving events...\n\n")

	for _, msg := range m.messages {
		if msg == nil {
			s += Nil.String() + "\n"
			continue
		}

		s += msg.String() + "\n"
	}

	if !m.quit {
		s += GrayForeground("\n\nPress any key to exit")
	}

	if m.quit {
		s += "\n"
	}

	return lipgloss.NewStyle().Margin(1, 2, 0, 2).Render(s)
}
