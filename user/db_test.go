package user

import (
	"encoding/json"
	"os"
	"testing"

	bolt "github.com/coreos/bbolt"
	"github.com/stretchr/testify/assert"
)

func setUpAccountDB() (*AccountDB, error) {
	return NewAccountDB("test")
}

func tearDownAccountDB(adb *AccountDB) {
	adb.db.Close()
	os.Remove("test.db")
}

func TestNewAccountDB(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)

	acc := Account{Name: "tealeg"}
	err = adb.Create(acc)
	assert.NoError(t, err)

	adb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("account"))
		key, _ := adb.Key(acc.Name)
		v := b.Get(key)
		acc2 := &Account{}
		err := json.Unmarshal(v, acc2)
		assert.NoError(t, err)
		assert.Equal(t, acc.Name, acc2.Name)
		return nil
	})
}

func TestCreateDuplicate(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)

	acc := Account{Name: "tealeg"}
	err = adb.Create(acc)
	assert.NoError(t, err)

	err = adb.Create(acc)
	assert.Error(t, err)

}

func TestGet(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)

	acc := Account{Name: "tealeg"}
	err = adb.Create(acc)
	assert.NoError(t, err)
	acc2, err := adb.Get("tealeg")
	assert.NoError(t, err)
	assert.Equal(t, acc.Name, acc2.Name)
}

func TestUpdate(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)

	acc1 := Account{Name: "tealeg", Forename: "Geoff"}
	err = adb.Create(acc1)
	assert.NoError(t, err)

	acc2 := Account{Name: "tealeg", Forename: "Gott"}
	err = adb.Update("tealeg", acc2)
	assert.NoError(t, err)

	acc3, err := adb.Get("tealeg")
	assert.NoError(t, err)
	assert.Equal(t, acc3.Forename, acc2.Forename)
	assert.Equal(t, acc3.Forename, acc2.Forename)
}

func TestAdminUserExists(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)

	exists, err := adb.AdminUserExists()
	assert.NoError(t, err)
	assert.False(t, exists)

	normalAcc := Account{Name: "nic", IsAdmin: false}
	err = adb.Create(normalAcc)
	assert.NoError(t, err)

	exists, err = adb.AdminUserExists()
	assert.NoError(t, err)
	assert.False(t, exists)

	adminAcc := Account{Name: "geoff", IsAdmin: true}
	err = adb.Create(adminAcc)
	assert.NoError(t, err)

	exists, err = adb.AdminUserExists()
	assert.NoError(t, err)
	assert.True(t, exists)
}
