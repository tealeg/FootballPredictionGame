package app

import (
	"github.com/labstack/echo"
	"github.com/tealeg/FPG2/competition"
	"github.com/tealeg/FPG2/user"
)

func Serve(port string, adb *user.AccountDB, cdb *competition.DB) {
	e := echo.New()

	setupUserHandlers(e, adb)
	setupFrontPageHandler(e, adb, cdb)
	setupLeagueHandlers(e, adb, cdb)
	e.Logger.Fatal(e.Start(port))
}
