package app

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/tealeg/FPG2/user"
)

func makeFrontPageHandler(adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		page := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>
    <h1>Welcome</h1>
  </body>
</html>
`
		return c.HTML(http.StatusOK, page)
	}
}

func setupFrontPageHandler(e *echo.Echo, adb *user.AccountDB) {
	e.GET("/frontpage", SecurePage(e, adb, makeFrontPageHandler(adb)))
}
