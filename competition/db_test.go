package league

import (
	"os"
	"testing"
)

func setUpDB(t *testing.T) *DB {
	db, err := NewDB("testDB")
	if err != nil {
		t.Fatalf("Unexpected error creating DB: %s", err.Error())
	}
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
