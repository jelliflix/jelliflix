package store

import (
	"github.com/go-redis/redis/v8"
	"github.com/jelliflix/jelliflix/infrastructure/datastore"
)

type Store struct {
	redis *redis.Client
}

func NewStore() *Store {
	r := datastore.NewRedis()
	return &Store{redis: r}
}

func (s *Store) AddMovie(imdbID string) (exists bool) {
	panic("implement me")
}

func (s *Store) SaveMovie(imdbID, filename string) {
	panic("implement me")
}

func (s *Store) RemoveMovie(imdbID string) {
	panic("implement me")
}

func (s *Store) AddEpisode(imdbID string) (exists bool) {
	panic("implement me")
}

func (s *Store) SaveEpisode(imdbID, filename string) {
	panic("implement me")
}

func (s *Store) RemoveEpisode(imdbID string) {
	panic("implement me")
}

func (s *Store) PurgeMovies(movies []string) {
	panic("implement me")
}

func (s *Store) PurgeSeries(episodes []string) {
	panic("implement me")
}
