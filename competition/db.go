package competition

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	bolt "github.com/coreos/bbolt"
)

type DB struct {
	db *bolt.DB
}

var (
	leagueBN []byte = []byte("league")
)

func NewDB(name string) (*DB, error) {
	db, err := bolt.Open(fmt.Sprintf("%s.db", name), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(leagueBN); err != nil {
			return err
		}

		return nil
	})
	return &DB{db: db}, nil
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

func (db *DB) Close() error {
	return db.db.Close()
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
