package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/jelliflix/jelliflix/internal/config"
)

// Bubble represents the properties of the UI.
type Bubble struct {
	keys     KeyMap
	config   config.Config
	spinner  spinner.Model
	messages []Message
}

// New creates a new instance of the UI.
func New() Bubble {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	b := Bubble{
		messages: make([]Message, 10),
		config:   cfg,
		keys:     DefaultKeyMap(),
	}

	b.spinner = spinner.New()
	b.spinner.Spinner = spinner.Dot
	b.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return b
}
