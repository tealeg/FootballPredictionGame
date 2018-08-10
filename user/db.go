package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bolt "github.com/coreos/bbolt"
)

type AccountDB struct {
	db         *bolt.DB
	bucketName []byte
}

func NewAccountDB(name string) (*AccountDB, error) {
	db, err := bolt.Open(fmt.Sprintf("%s.db", name), 0600,
		&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	bname := []byte("account")
	err = db.Update(func(tx *bolt.Tx) error {
		var err error

		_, err = tx.CreateBucketIfNotExists(bname)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	adb := &AccountDB{
		db:         db,
		bucketName: bname,
	}
	return adb, nil
}

// Create stores a Account in the account DB
func (adb *AccountDB) Create(account Account) error {
	return adb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(adb.bucketName)
		buf, err := json.Marshal(account)
		if err != nil {
			return err
		}
		key, err := adb.Key(account.Name)
		if err != nil {
			return err
		}
		v := b.Get(key)
		if len(v) != 0 {
			return errors.New("Duplicate key")
		}
		return b.Put(key, buf)
	})
}

// Get retrieves a user Account with the provided name
func (adb *AccountDB) Get(name string) (*Account, error) {
	acc := &Account{}
	err := adb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(adb.bucketName)
		key, err := adb.Key(name)
		if err != nil {
			return err
		}
		v := b.Get(key)
		if len(v) == 0 {
			return errors.New("Zero bytes returned for key")
		}
		err = json.Unmarshal(v, acc)
		if err != nil {
			return err
		}
		return nil
	})
	return acc, err
}

func (adb *AccountDB) Update(name string, account Account) error {
	return adb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(adb.bucketName)
		buf, err := json.Marshal(account)
		if err != nil {
			return err
		}
		key, err := adb.Key(account.Name)
		if err != nil {
			return err
		}
		return b.Put(key, buf)
	})

}

//
func (adb *AccountDB) AdminUserExists() (bool, error) {
	found := false
	err := adb.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(adb.bucketName)

		c := b.Cursor()

		acc := &Account{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, acc)
			if err != nil {
				return err
			}
			if acc.IsAdmin {
				found = true
				return nil
			}
		}

		return nil
	})
	return found, err
}

// Generate a key from an account name
func (adb *AccountDB) Key(name string) ([]byte, error) {
	if name == "" {
		return nil, errors.New("No Name set for account")
	}
	return []byte(name), nil
}

func (adb *AccountDB) Close() error {
	return adb.db.Close()
}
