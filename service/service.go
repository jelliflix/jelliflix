package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/jelliflix/imdb"
	"github.com/jelliflix/jelliflix/store"
	"github.com/jelliflix/jelliflix/ui"
	"github.com/jelliflix/meta"
	finder "github.com/jelliflix/torrent"
)

var (
	timeout         = time.Second * 10
	moviesQuality   = "1080p"
	episodesQuality = "720p"
)

type Service struct {
	path    string
	imdb    *imdb.IMDB
	store   *store.Store
	finder  *finder.Torrent
	torrent *torrent.Client
	program *tea.Program
}

func NewService(program *tea.Program, id, path string, store *store.Store) *Service {
	i, err := imdb.NewIMDB(imdb.DefaultOptions, id)
	if err != nil {
		program.Send(ui.ErrorMsg{Error: err})
	}

	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = path
	cfg.ListenPort = 42067
	t, err := torrent.NewClient(cfg)
	if err != nil {
		program.Send(ui.ErrorMsg{Error: err})
	}

	cache := finder.NewInMemCache()
	cinemeta := meta.NewCinemeta(meta.DefaultOptions)
	clients := []finder.Client{
		finder.NewYTS(finder.DefaultYTSOpts, cache),
		finder.NewRARBG(finder.DefaultRARBOpts, cache),
		finder.NewTPB(finder.DefaultTPBOpts, cache, cinemeta),
	}

	f := finder.NewTorrent(clients, timeout)

	return &Service{
		imdb:    i,
		finder:  f,
		torrent: t,
		path:    path,
		store:   store,
		program: program,
	}
}

func (s *Service) ProcessWatchList(ctx context.Context) {
	for ; true; <-time.NewTicker(timeout).C {
		s.program.Send(ui.InfoMsg{Message: ui.Bold("EXPORTING WATCHLIST [IMDb]...")})
		watchList, err := s.imdb.ExportWatchList()
		if err != nil {
			s.program.Send(ui.ErrorMsg{Error: err})
		}
		if len(watchList) == 0 {
			continue
		}
		s.program.Send(ui.InfoMsg{Message: ui.Bold(ui.RedForeground("PURGE WATCHED [IMDb]..."))})
		// TODO: Purge watched titles
		mustPurge := s.store.Items()

		for _, item := range watchList {
			delete(mustPurge, item.ID)
			s.findTorrents(ctx, item)
		}
	}
}

func (s *Service) findTorrents(ctx context.Context, item imdb.Item) {
	if s.store.Exists(item.ID) {
		return
	}

	switch item.Type {
	case "movie":
		results, err := s.finder.FindMovie(ctx, item.ID)
		if err != nil {
			s.program.Send(ui.ErrorMsg{Error: err})
		}
		s.selectTorrent(item, results, moviesQuality)
	case "tvEpisode":
		results, err := s.finder.FindEpisode(ctx, item.ID, item.Name)
		if err != nil {
			s.program.Send(ui.ErrorMsg{Error: err})
		}
		s.selectTorrent(item, results, episodesQuality)
	}
}

func (s *Service) selectTorrent(item imdb.Item, results []finder.Result, quality string) {
	var (
		sameQuality  []finder.Result
		downloadable finder.Result
	)
	for _, result := range results {
		if result.Quality == quality {
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

	s.downloadTorrent(item, downloadable)
}

func (s *Service) downloadTorrent(item imdb.Item, result finder.Result) {
	t, err := s.torrent.AddMagnet(result.MagnetURL)
	if err != nil {
		s.program.Send(ui.ErrorMsg{Error: err})
	}

	<-t.GotInfo()
	s.program.Send(ui.InfoMsg{
		Message: fmt.Sprintf("DOWNLOADING [%s %s] %s",
			ui.Italic(ui.PrimaryForeground(result.Quality)),
			ui.Italic(ui.PrimaryForeground(humanize.Bytes(uint64(result.Size)))),
			ui.Bold(ui.SecondaryForeground(strings.ToUpper(fmt.Sprintf("%s (%d)", item.Name, item.Year)))),
		),
	})
	if err = s.store.Set(item.ID, t.Name()); err != nil {
		s.program.Send(ui.ErrorMsg{Error: err})
	}

	go t.DownloadAll()
	// TODO: Announce download progress via integration
	// announceProgress(t, result)

	return
}

func (s *Service) announceProgress(t *torrent.Torrent, r finder.Result) {
	for {
		s.program.Send(ui.InfoMsg{
			Message: fmt.Sprintf("DOWNLOADED %s/%s",
				ui.Italic(ui.PrimaryForeground(humanize.Bytes(uint64(t.BytesCompleted())))),
				ui.Italic(ui.PrimaryForeground(humanize.Bytes(uint64(r.Size)))),
			),
		})
		time.Sleep(time.Second * 5)

		if t.Complete.Bool() {
			break
		}
	}
}
