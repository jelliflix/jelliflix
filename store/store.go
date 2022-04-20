package store

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jelliflix/jelliflix/infrastructure/datastore"
)

var (
	Movies   Kind = "M"
	Episodes Kind = "EP"
)

type Kind string

type Store struct {
	redis *redis.Client
}

func NewStore() *Store {
	r := datastore.NewRedis()
	return &Store{redis: r}
}

func (s *Store) Set(ctx context.Context, kind Kind, imdbID, filename string) {
	s.redis.Set(ctx, s.getKey(kind, imdbID), filename, 0)
}

func (s *Store) Get(ctx context.Context, kind Kind) (items map[string]string) {
	keys := s.redis.Keys(ctx, s.getKey(kind, "*")).Val()

	for _, key := range keys {
		filename := s.redis.Get(ctx, key).Val()
		items[key] = filename
	}

	return
}

func (s *Store) Purge(ctx context.Context, kind Kind, items []string) {
	panic("implement me")
}

func (s *Store) Remove(ctx context.Context, kind Kind, imdbID string) {
	s.redis.Del(ctx, s.getKey(kind, imdbID))
}

func (s *Store) Exists(ctx context.Context, kind Kind, imdbID string) bool {
	return s.redis.Get(ctx, s.getKey(kind, imdbID)).Val() != ""
}

func (s *Store) getKey(kind Kind, imdbID string) string {
	switch kind {
	case Movies:
		return s.moviesKey(imdbID)
	case Episodes:
		return s.episodesKey(imdbID)
	}

	return ""
}

func (s *Store) moviesKey(imdbID string) string {
	return fmt.Sprintf("m:%s", imdbID)
}

func (s *Store) episodesKey(imdbID string) string {
	return fmt.Sprintf("ep:%s", imdbID)
}
