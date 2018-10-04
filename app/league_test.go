package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/tealeg/FootballPredictionGame/competition"
)

// newLeagueRequest.Validate indicates all validation failures.
func TestNewLeagueRequestValidation(t *testing.T) {
	cases := []struct {
		NLR         *newLeagueRequest
		Expectation []string
	}{
		{
			NLR:         &newLeagueRequest{},
			Expectation: []string{"League name is empty"},
		},
		{
			NLR:         &newLeagueRequest{Name: "English Premier League"},
			Expectation: []string{},
		},
	}

	for i, c := range cases {
		eLen := len(c.Expectation)
		r := newSimpleResponse()
		err := c.NLR.Validate(r)
		if eLen > 0 {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		for j, e := range c.Expectation {
			er := r.Errors[j]
			assert.Equalf(t, e, er, "Case %d, Expectation %d: validation error == %q, should be %q", i, j, e, er)
		}
	}
}

// newLeagueRequest.createLeaugue creates a competition.League object and persists it.
func TestNewLeagueRequestCreateLeague(t *testing.T) {
	cdb, err := setUpCompetitionDB()
	assert.NoError(t, err)
	defer tearDownCompetitionDB(cdb)

	nlr := newLeagueRequest{Name: "English Premier League"}
	id, err := nlr.createLeague(cdb)
	assert.NoError(t, err)

	l, err := cdb.GetLeague(id)
	assert.NoError(t, err)

	assert.Equal(t, nlr.Name, l.Name)
}

// GetAllLeagues list all the leagues in the db.
func TestGetAllLeaguesHandler(t *testing.T) {
	cdb, err := setUpCompetitionDB()
	assert.NoError(t, err)
	defer tearDownCompetitionDB(cdb)

	var expected []competition.League
	for i := 0; i < 5; i++ {
		l := &competition.League{
			Name: fmt.Sprintf("League %d", i+1),
		}
		_, err := cdb.CreateLeague(l)
		assert.NoError(t, err)
		expected = append(expected, *l)
	}
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/leagues.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeGetAllLeaguesHandler(e, cdb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body)

	var result []competition.League
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 5)

	for i, l := range result {
		exp := expected[i]
		assert.Equal(t, exp.ID, l.ID)
		assert.Equal(t, exp.Name, l.Name)
	}
}

func TestNewLeagueHandler(t *testing.T) {
	cdb, err := setUpCompetitionDB()
	assert.NoError(t, err)
	defer tearDownCompetitionDB(cdb)
	
	e := echo.New()
	expectedName := "English Premier League"
	encoded := fmt.Sprintf(`{"name": %q}`, expectedName)

	req := httptest.NewRequest(echo.POST, "/leagues/new.json", strings.NewReader(encoded))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := makeNewLeagueHandler(e, cdb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	leagues, err := cdb.GetAllLeagues()
	assert.NoError(t, err)

	assert.Len(t, leagues, 1)

	l := leagues[0]
	assert.Equal(t, expectedName, l.Name)
}

func TestLeagueHandler(t *testing.T) {
	cdb, err := setUpCompetitionDB()
	assert.NoError(t, err)
	defer tearDownCompetitionDB(cdb)

	l := competition.League{Name: "1. Bundesliga"}
	lid, err := cdb.CreateLeague(&l)
	assert.NoError(t, err)

	leagueID := strconv.FormatUint(lid, 10)
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/league/"+leagueID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/league/:id")
	c.SetParamNames("id")
	c.SetParamValues(leagueID)

	h := makeLeagueHandler(e, cdb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body.Bytes())

	lr := LeagueResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), &lr)
	assert.NoError(t, err)

	assert.Equal(t, l.ID, lr.ID)
	assert.Equal(t, l.Name, lr.Name)
	assert.Empty(t, lr.Seasons)
}

func TestLeagueHandlerWithSeasons(t *testing.T) {
	cdb, err := setUpCompetitionDB()
	assert.NoError(t, err)
	defer tearDownCompetitionDB(cdb)

	l := competition.League{Name: "1. Bundesliga"}
	lid, err := cdb.CreateLeague(&l)
	assert.NoError(t, err)

	s := competition.NewSeason(lid, 2018, 2019)
	_, err = cdb.CreateSeason(s)
	assert.NoError(t, err)

	leagueID := strconv.FormatUint(lid, 10)
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/league/"+leagueID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/league/:id")
	c.SetParamNames("id")
	c.SetParamValues(leagueID)

	h := makeLeagueHandler(e, cdb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body.Bytes())

	lr := LeagueResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), &lr)
	assert.NoError(t, err)

	assert.Equal(t, l.ID, lr.ID)
	assert.Equal(t, l.Name, lr.Name)
	assert.Len(t, lr.Seasons, 1)

	sr := lr.Seasons[0]
	assert.Equal(t, s.ID, sr.ID)
	assert.Equal(t, s.StartYear, sr.StartYear)
	assert.Equal(t, s.EndYear, sr.EndYear)
}
