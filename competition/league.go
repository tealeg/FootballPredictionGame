package competition

import (
	"encoding/json"

	bolt "github.com/coreos/bbolt"
)

type League struct {
	ID   uint64
	Name string
}

func (db *DB) CreateLeague(l *League) (uint64, error) {
	var id uint64

	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(leagueBN)
		id, _ = b.NextSequence()
		l.ID = id

		// Marshal user data into bytes.
		buf, err := json.Marshal(l)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(l.ID), buf)
	})
	return id, err
}

func (db *DB) GetLeague(id uint64) (*League, error) {
	l := &League{}
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(leagueBN)
		v := b.Get(itob(id))
		err := json.Unmarshal(v, l)
		if err != nil {
			return err
		}
		return nil
	})
	return l, err
}

func (db *DB) GetAllLeagues() ([]League, error) {
	leagues := make([]League, 0)
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(leagueBN)
		err := b.ForEach(func(k, v []byte) error {
			l := League{}
			err := json.Unmarshal(v, &l)
			if err != nil {
				return err
			}
			leagues = append(leagues, l)
			return nil
		})
		return err

	})
	return leagues, err
}
