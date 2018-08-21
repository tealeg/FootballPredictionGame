package main

import (
	"log"

	"github.com/tealeg/FootballPredictionGame/app"
	"github.com/tealeg/FootballPredictionGame/competition"
	"github.com/tealeg/FootballPredictionGame/user"
)

func main() {
	adb, err := user.NewAccountDB("user")
	if err != nil {
		log.Fatal(err)
	}
	cdb, err := competition.NewDB("competition")
	if err != nil {
		log.Fatal(err)
	}

	app.Serve(":9090", adb, cdb)
}
