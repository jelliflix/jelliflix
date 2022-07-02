package exporter

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jelliflix/imdb"
	"github.com/jelliflix/jelliflix/internal/config"
	"github.com/jelliflix/jelliflix/internal/storage"
	"github.com/jelliflix/jelliflix/internal/ui"
	"github.com/jelliflix/torrent"
)

type Exporter struct {
	program *tea.Program
	config  config.Config

	imdb    *imdb.IMDB
	torrent *torrent.Torrent
	storage *storage.Storage
}

func New(cfg config.Config, p *tea.Program) Exporter {
	return Exporter{config: cfg, program: p}
}

func (e Exporter) Start() error {
	var err error
	e.imdb, err = imdb.NewIMDB(imdb.DefaultOptions, e.config.Providers.IMDb.ID)
	if err != nil {
		return err
	}

	d, err := time.ParseDuration(e.config.Settings.RefreshRate)
	if err != nil {
		return err
	}

	go func() {
		for ; true; <-time.NewTicker(d).C {
			watchList, err := e.imdb.ExportWatchList()
			if err != nil {
				e.program.Send(ui.ErrorMsg{Error: err})
			}
			for _, item := range watchList {
				e.program.Send(ui.InfoMsg{Message: item.ID})
			}
		}
	}()

	return err
}
