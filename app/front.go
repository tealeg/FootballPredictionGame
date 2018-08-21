package app

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tealeg/FPG2/competition"
	"github.com/tealeg/FPG2/user"
)

func makeFrontPageHandler(adb *user.AccountDB, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		pageHead := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>
    <h1>Welcome</h1>`

		pageFoot := `
  </body>
</html>
`
		leagues := "<h2>Leagues</h2>"
		leagues += "<ul>"
		all, err := cdb.GetAllLeagues()
		if err != nil {
			leagues += fmt.Sprintf("<li class=\"error\">Error getting leagues: %s</li>", err.Error())
		}
		for _, l := range all {
			leagues += fmt.Sprintf("<li><a href=\"/league/%d\">%s</a></li>", l.ID, l.Name)
		}
		leagues += "</ul>"
		user, err := GetUserAccount(c, adb)
		if err == nil {
			if user.IsAdmin {
				leagues += "<form action=\"/league/new\" method=\"POST\"><fieldset>League Name: <input type=\"text\" name=\"name\" /><br /><input type=\"submit\" value=\"Create\"/></fieldset></form>"
			}
		}
		page := pageHead + leagues + pageFoot
		return c.HTML(http.StatusOK, page)
	}
}

func setupFrontPageHandler(e *echo.Echo, adb *user.AccountDB, cdb *competition.DB) {
	e.GET("/frontpage", SecurePage(e, adb, makeFrontPageHandler(adb, cdb)))
}
