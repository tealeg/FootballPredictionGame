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

func makeLeagueHandler(e *echo.Echo, adb *user.AccountDB, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		l, err := cdb.GetLeague(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		seasons, err := cdb.GetAllLeagueSeasons(l.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		sSection := "<ul>"
		for _, s := range seasons {
			sSection += fmt.Sprintf("<li><a href=\"/season/%s\">%d/%d</a></li>", s.ID, s.StartYear, s.EndYear)
		}
		sSection += "</ul>"

		createForm := ""
		user, err := GetUserAccount(c, adb)
		if err == nil {
			if user.IsAdmin {
				createForm += fmt.Sprintf(`
<form action="/season/new" method="POST">
  <fieldset>
    <input type="hidden" name="leagueID" value="%d" />
    Season Start Year: <input type="text" name="startYear" placeholder="YYYY"/><br />
    Season End Year: <input type="text" name="endYear" placeholder="YYYY"/>
    <br />
    <input type="submit" value="Create Season"/>
  </fieldset>
</form>`, l.ID)
			}
		}

		return c.HTML(http.StatusOK, fmt.Sprintf(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>
    <h1>%s</h1>
    %s
    %s
  </body>
</html>
`, l.Name, sSection, createForm))
	}
}

func setupLeagueHandlers(e *echo.Echo, adb *user.AccountDB, cdb *competition.DB) {
	e.POST("/league/new", SecurePage(e, adb, makeNewLeagueHandler(e, cdb)))
	e.GET("/league/:id", SecurePage(e, adb, makeLeagueHandler(e, adb, cdb)))
}
