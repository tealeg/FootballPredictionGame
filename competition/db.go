package competition

import (
	"encoding/binary"
	"fmt"
	"time"

	bolt "github.com/coreos/bbolt"
)

type DB struct {
	db *bolt.DB
}

var (
	leagueBN []byte   = []byte("league")
	seasonBN []byte   = []byte("season")
	BNs      [][]byte = [][]byte{
		leagueBN,
		seasonBN,
	}
)

func NewDB(name string) (*DB, error) {
	db, err := bolt.Open(fmt.Sprintf("%s.db", name), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bn := range BNs {
			if _, err := tx.CreateBucketIfNotExists(bn); err != nil {
				return err
			}

		}

		return nil
	})
	return &DB{db: db}, nil
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
