package exporter

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/jelliflix/imdb"
	"github.com/jelliflix/jelliflix/internal/config"
	"github.com/jelliflix/jelliflix/internal/storage"
	"github.com/jelliflix/jelliflix/internal/ui"
	"github.com/jelliflix/meta"
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

func (e Exporter) Start(ctx context.Context) error {
	var err error
	var clients []torrent.Client
	cache := torrent.NewInMemCache()
	cinemeta := meta.NewCinemeta(meta.DefaultOptions)
	timeout, err := time.ParseDuration(e.config.Settings.Timeout)
	if err != nil {
		return err
	}

	for _, provider := range e.config.Exporter.Providers {
		switch strings.ToUpper(provider) {
		case "YTS":
			clients = append(clients, torrent.NewYTS(torrent.DefaultYTSOpts, cache))
		case "TPB":
			clients = append(clients, torrent.NewTPB(torrent.DefaultTPBOpts, cache, cinemeta))
		case "RARBG":
			clients = append(clients, torrent.NewRARBG(torrent.DefaultRARBOpts, cache))
		}
	}

	if len(clients) == 0 {
		return fmt.Errorf("missing torrent providers")
	}

	e.torrent = torrent.NewTorrent(clients, timeout)
	e.storage = storage.New(e.config.Settings.StorageDataDir)
	e.imdb, err = imdb.NewIMDB(imdb.DefaultOptions, e.config.Providers.IMDb.ID)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(e.config.Settings.RefreshRate)
	if err != nil {
		return err
	}

	go func() {
		for ; true; <-time.NewTicker(duration).C {
			e.program.Send(ui.InfoMsg{Message: "EXPORTING WATCHLIST [IMDb]..."})
			watchList, err := e.imdb.ExportWatchList()
			if err != nil {
				e.program.Send(ui.ErrorMsg{Error: err})
			}

			for _, item := range watchList {
				e.findTorrents(ctx, item)
			}
		}
	}()

	return err
}

func (e Exporter) findTorrents(ctx context.Context, item imdb.Item) {
	if e.storage.Exists(item.ID) {
		return
	}

	var err error
	var results []torrent.Result
	switch item.Type {
	case "movie":
		results, err = e.torrent.FindMovie(ctx, item.ID)
		if err != nil {
			e.program.Send(ui.ErrorMsg{Error: err})
		}
	case "tvEpisode":
		results, err = e.torrent.FindEpisode(ctx, item.ID, item.Name)
		if err != nil {
			e.program.Send(ui.ErrorMsg{Error: err})
		}
	}

	title := strings.ToUpper(fmt.Sprintf("%s (%d)", item.Name, item.Year))

	if len(results) == 0 {
		// TODO -> Add titles with no results to retry queue
		e.program.Send(ui.InfoMsg{Message: fmt.Sprintf("NO RESULTS FOR %s", title)})
	} else {
		e.program.Send(ui.InfoMsg{Message: fmt.Sprintf("FOUND %d RESULTS FOR %s", len(results), title)})
		e.selectTorrent(title, results)
	}
}

func (e Exporter) selectTorrent(title string, results []torrent.Result) {
	var downloadable torrent.Result
	var sameQuality []torrent.Result
	for _, result := range results {
		if result.Quality == e.config.Settings.DownloadQuality {
			sameQuality = append(sameQuality, result)
		}
	}

	if len(sameQuality) == 0 {
		downloadable = results[0]
	} else {
		sort.Slice(sameQuality, func(i, j int) bool {
			return sameQuality[i].Size < sameQuality[j].Size
		})
		downloadable = sameQuality[0]
	}

	// TODO -> Download the title with torrent package (needs implementation)
	e.program.Send(ui.InfoMsg{Message: fmt.Sprintf("STARTED DOWNLOADING [%s %s] %s",
		downloadable.Quality, humanize.Bytes(uint64(downloadable.Size)), title),
	})
}
