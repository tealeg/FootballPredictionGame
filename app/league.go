package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tealeg/FPG2/competition"
	"github.com/tealeg/FPG2/user"
)

func makeNewLeagueHandler(e *echo.Echo, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		if name == "" {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid league name: %q", name))
		}
		league := &competition.League{Name: name}
		id, err := cdb.CreateLeague(league)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		url := fmt.Sprintf("/league/%d", id)
		return c.Redirect(http.StatusSeeOther, url)
	}
}

func makeLeagueHandler(e *echo.Echo, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		l, err := cdb.GetLeague(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.HTML(http.StatusOK, fmt.Sprintf(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>
    <h1>%s</h1>
  </body>
</html>
`, l.Name))
	}
}

func setupLeagueHandlers(e *echo.Echo, adb *user.AccountDB, cdb *competition.DB) {
	e.POST("/league/new", SecurePage(e, adb, makeNewLeagueHandler(e, cdb)))
	e.GET("/league/:id", SecurePage(e, adb, makeLeagueHandler(e, cdb)))
}
