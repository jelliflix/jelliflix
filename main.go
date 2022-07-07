package main

import (
	"context"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jelliflix/jelliflix/internal/config"
	"github.com/jelliflix/jelliflix/internal/exporter"
	"github.com/jelliflix/jelliflix/internal/ui"
	"github.com/urfave/cli/v2"
)

const Version = "undefined"

func main() {
	app := cli.NewApp()
	app.Name = "jelliflix"
	app.Usage = "Organize video content from IMDB"
	app.Version = Version

	app.Action = func(_ *cli.Context) error {
		cfg, err := config.ParseConfig()
		if err != nil {
			log.Fatal(err)
		}

		// If logging is enabled, logs will be output to debug.log.
		if cfg.Settings.EnableLogging {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				log.Fatal(err)
			}

			defer func() {
				if err = f.Close(); err != nil {
					log.Fatal(err)
				}
			}()
		}

		m := ui.New()
		var opts []tea.ProgramOption

		// Always append alt screen program option.
		opts = append(opts, tea.WithAltScreen())

		// Initialize and start app.
		p := tea.NewProgram(m, opts...)
		e := exporter.New(cfg, p)
		if err = e.Start(context.Background()); err != nil {
			log.Fatalf("failed to start exporter: %v", err)
		}
		if err = p.Start(); err != nil {
			log.Fatalf("failed to start jelliflix: %v", err)
		}

		return err
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
