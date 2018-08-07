package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bolt "github.com/coreos/bbolt"
)

type AccountDB struct {
	db *bolt.DB
}

func NewAccountDB(name string) (*AccountDB, error) {
	db, err := bolt.Open(fmt.Sprintf("%s.db", name), 0600,
		&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	var b *bolt.Bucket
	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		b, err = tx.CreateBucket([]byte("account"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	adb := &AccountDB{
		db: db,
	}
	return adb, nil
}

// Create stores a Account in the account DB
func (adb *AccountDB) Create(account Account) error {
	return adb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("account"))
		buf, err := json.Marshal(account)
		if err != nil {
			return err
		}
		key, err := adb.Key(account)
		if err != nil {
			return err
		}
		return b.Put(key, buf)
	})
}

// Generate a key for the account bucket
func (adb *AccountDB) Key(account Account) ([]byte, error) {
	if account.Name == "" {
		return nil, errors.New("No Name set for account")
	}
	return []byte(account.Name), nil
}
