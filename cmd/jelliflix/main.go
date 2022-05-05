package main

import (
	"context"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jelliflix/jelliflix/service"
	"github.com/jelliflix/jelliflix/store"
	"github.com/jelliflix/jelliflix/ui"
	"github.com/urfave/cli/v2"
)

const version = "v0.0.1-beta"

func main() {
	app := cli.NewApp()
	app.Name = "jelliflix"
	app.Usage = "Store and organize video content from IMDB."
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "imdb",
			Usage:   "IMDb id",
			Aliases: []string{"i"},
		},
		&cli.StringFlag{
			Name:    "logs",
			Usage:   "logs path",
			Aliases: []string{"l"},
		},
		&cli.StringFlag{
			Name:    "path",
			Usage:   "download path",
			Aliases: []string{"p"},
		},
	}
	app.Version = version
	app.Action = func(c *cli.Context) (err error) {
		logs := c.String("logs")
		if logs != "" {
			log.SetFlags(0)
			f, err := tea.LogToFile(logs, "jelliflix")
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
			defer func() { _ = f.Close() }()
		}

		imdb := c.String("imdb")
		if imdb == "" {
			return cli.Exit("missing imdb user/list id argument", 1)
		}

		path := c.String("path")
		if path == "" {
			return cli.Exit("missing jelliflix path argument", 1)
		}

		p := tea.NewProgram(ui.NewApp(), tea.WithAltScreen())
		s := service.NewService(p, imdb, path, store.NewStore(path))

		go s.ProcessWatchList(context.Background())
		if err := p.Start(); err != nil {
			return cli.Exit(err.Error(), 1)
		}

		return
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
