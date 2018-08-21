package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tealeg/FPG2/competition"
	"github.com/tealeg/FPG2/user"
)

func makeNewSeasonHandler(e *echo.Echo, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		leagueID, err := strconv.Atoi(c.FormValue("leagueID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		startYear, err := strconv.Atoi(c.FormValue("startYear"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		endYear, err := strconv.Atoi(c.FormValue("endYear"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		s := competition.NewSeason(uint64(leagueID), uint16(startYear), uint16(endYear))
		id, err := cdb.CreateSeason(s)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		url := fmt.Sprintf("/season/%s", id)
		return c.Redirect(http.StatusSeeOther, url)
	}
}

func makeSeasonHandler(e *echo.Echo, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		s, err := cdb.GetSeason(id)
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
    <h1>Season %d/%d</h1>
  </body>
</html>`, s.StartYear, s.EndYear))
	}
}

func setupSeasonHandlers(e *echo.Echo, adb *user.AccountDB, cdb *competition.DB) {
	e.POST("/season/new", SecurePage(e, adb, makeNewSeasonHandler(e, cdb)))
	e.GET("/season/:id", SecurePage(e, adb, makeSeasonHandler(e, cdb)))
}
