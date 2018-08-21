package app

import (
	"github.com/labstack/echo"
	"github.com/tealeg/FootballPredictionGame/competition"
	"github.com/tealeg/FootballPredictionGame/user"
)

func Serve(port string, adb *user.AccountDB, cdb *competition.DB) {
	e := echo.New()
	e.Static("/", "static")
	setupUserHandlers(e, adb)
	setupFrontPageHandler(e, adb, cdb)
	setupLeagueHandlers(e, adb, cdb)
	setupSeasonHandlers(e, adb, cdb)
	e.Logger.Fatal(e.Start(port))
}
