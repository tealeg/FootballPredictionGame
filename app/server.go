package app

import (
	"github.com/labstack/echo"
	"github.com/tealeg/FPG2/user"
)

func Serve(port string, adb *user.AccountDB) {
	e := echo.New()

	setupUserHandlers(e, adb)
	setupFrontPageHandler(e, adb)
	e.Logger.Fatal(e.Start(port))
}
