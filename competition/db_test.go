package competition

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setUpDB(t *testing.T) *DB {
	db, err := NewDB("testDB")
	assert.NoError(t, err)
	return db
}

func tearDownDB(db *DB) {
	db.Close()
	os.Remove("testDB.db")

}

func TestNewDB(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(db)
}
