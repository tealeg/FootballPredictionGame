package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tealeg/FootballPredictionGame/competition"
	"github.com/tealeg/FootballPredictionGame/user"
)

// newLeagueRequest is a holder for data passed into new League requests.
type newLeagueRequest struct {
	Name string
}

// Validate checks the members of a newLeagueRequest for validity and
// populates a simpleResponse with the errors it finds.  The last
// error found will be returned, and can be used to indicate overall
// validation failure (or, if nil, success).
func (nlr *newLeagueRequest) Validate(r *simpleResponse) error {
	var err error
	if nlr.Name == "" {
		err = errors.New("League name is empty")
		r.AddError(err)
	}
	return err
}

//createLeague creates a competition.League based on its newLeagueRequest.
func (nlr *newLeagueRequest) createLeague(cdb *competition.DB) (uint64, error) {
	league := &competition.League{Name: nlr.Name}
	return cdb.CreateLeague(league)
}

func makeNewLeagueHandler(e *echo.Echo, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		e.Logger.Info("Creating League")
		nlr := new(newLeagueRequest)
		if err := c.Bind(nlr); err != nil {
			e.Logger.Error(err.Error())
			return err
		}

		r := newSimpleResponse()
		err := nlr.Validate(r)
		if err != nil {
			for _, rerr := range r.Errors {
				e.Logger.Error(rerr)
			}
			return c.JSON(http.StatusBadRequest, r.Errors)
		}

		lid, err := nlr.createLeague(cdb)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		r.ObjID = strconv.FormatUint(lid, 10)
		e.Logger.Info("User created: %+v", r)
		return c.JSON(http.StatusOK, r)
	}
}

func makeGetAllLeaguesHandler(e *echo.Echo, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ls, err := cdb.GetAllLeagues()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, &ls)
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
	e.GET("/leagues.json", SecurePage(e, adb, makeGetAllLeaguesHandler(e, cdb)))
	e.POST("/leagues/new.json", SecurePage(e, adb, makeNewLeagueHandler(e, cdb)))
	e.GET("/league/:id", SecurePage(e, adb, makeLeagueHandler(e, adb, cdb)))
}
