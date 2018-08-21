package app

import (
	"os"

	"github.com/tealeg/FootballPredictionGame/user"
)

func setUpAccountDB() (*user.AccountDB, error) {
	return user.NewAccountDB("test-account")
}

func tearDownAccountDB(adb *user.AccountDB) {
	adb.Close()
	os.Remove("test-account.db")
}
