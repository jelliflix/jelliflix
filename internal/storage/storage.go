package storage

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

type Storage struct {
	db *badger.DB
}

func New(path string) *Storage {
	c, err := badger.Open(
		badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR))
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{db: c}
}

func (s *Storage) Set(id, state string) (err error) {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	err = txn.Set([]byte(id), []byte(state))
	if err != nil {
		return
	}

	if err = txn.Commit(); err != nil {
		return
	}

	return
}

func (s *Storage) Get(id string) (state string) {
	_ = s.db.View(func(txn *badger.Txn) (err error) {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		_ = item.Value(func(val []byte) (err error) {
			state = string(val)

			return
		})

		return
	})

	return state
}

func (s *Storage) Items() (items map[string]string) {
	items = make(map[string]string)
	_ = s.db.View(func(txn *badger.Txn) (err error) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err = item.Value(func(v []byte) (err error) {
				items[string(k)] = string(v)
				return
			})
			if err != nil {
				return
			}
		}
		return
	})

	return
}

func (s *Storage) Remove(id string) {
	_ = s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(id)); err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (s *Storage) Exists(id string) (exists bool) {
	_ = s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(id))
		exists = err == nil
		return err
	})

	return
}
