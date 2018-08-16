package app

import (
	"os"

	"github.com/tealeg/FPG2/competition"
)

func setUpCompetitionDB() (*competition.DB, error) {
	return competition.NewDB("test-competition")
}

func tearDownCompetitionDB(cdb *competition.DB) {
	cdb.Close()
	os.Remove("test-competition.db")
}
