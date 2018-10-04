package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tealeg/FootballPredictionGame/competition"
	"github.com/tealeg/FootballPredictionGame/user"
)

// newLeagueRequest is a holder for data passed into new League requests.
type newLeagueRequest struct {
	Name string `json:"name"`
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
		e.Logger.Info("League created: %+v", r)
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

type SeasonResponse struct {
	ID        string `json:"id"`
	StartYear uint16 `json:"startyear"`
	EndYear   uint16 `json:"endyear"`
}

type LeagueResponse struct {
	ID      uint64           `json:"id"`
	Name    string           `json:"name"`
	Seasons []SeasonResponse `json:"seasons"`
}

func makeLeagueHandler(e *echo.Echo, cdb *competition.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sid := c.Param("id")
		id, err := strconv.ParseUint(sid, 10, 64)
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
		lr := LeagueResponse{
			ID:   l.ID,
			Name: l.Name,
		}
		for _, s := range seasons {
			sr := SeasonResponse{
				ID:        s.ID,
				StartYear: s.StartYear,
				EndYear:   s.EndYear,
			}
			lr.Seasons = append(lr.Seasons, sr)
		}

		return c.JSON(http.StatusOK, lr)
	}
}

func setupLeagueHandlers(e *echo.Echo, adb *user.AccountDB, cdb *competition.DB) {
	e.GET("/leagues.json", SecurePage(e, adb, makeGetAllLeaguesHandler(e, cdb)))
	e.POST("/leagues/new.json", SecurePage(e, adb, makeNewLeagueHandler(e, cdb)))
	e.GET("/league/:id", SecurePage(e, adb, makeLeagueHandler(e, cdb)))
}
