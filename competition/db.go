package league

import (
	"fmt"
	"time"

	bolt "github.com/coreos/bbolt"
)

type DB struct {
	db *bolt.DB
}

func NewDB(name string) (*DB, error) {
	db, err := bolt.Open(fmt.Sprintf("%s.db", name), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("league")); err != nil {
			return err
		}

		return nil
	})
	return &DB{db: db}, nil
}

//
func (db *DB) Close() error {
	return db.db.Close()
}
