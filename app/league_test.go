package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	if err != nil {
		t.Fatalf("unexpected error setting up competition db: %s", err.Error())
	}
	defer tearDownCompetitionDB(cdb)
	var expected []competition.League
	for i := 0; i < 5; i++ {
		l := &competition.League{
			Name: fmt.Sprintf("League %d", i+1),
		}
		_, err := cdb.CreateLeague(l)
		if err != nil {
			t.Fatalf("Error creating league %d: %s", i+1, err)
		}
		expected = append(expected, *l)
	}
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/leagues.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeGetAllLeaguesHandler(e, cdb)
	err = h(c)
	if err != nil {
		t.Fatalf("Unexepected error in handler: %s", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected rec.Code = OK, but got %s", http.StatusText(rec.Code))
	}
	if rec.Body.Len() == 0 {
		t.Error("Empty body")
	}
	var result []competition.League
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("error unmarshalling body: %s", err.Error())
	}
	if len(result) != 5 {
		t.Errorf("expected 5 leagues, but got %d", len(result))
	}
	for i, l := range result {
		exp := expected[i]
		if exp.ID != l.ID {
			t.Errorf("Result[%d]: expected ID == %d, but got %d", i, exp.ID, l.ID)
		}
		if exp.Name != l.Name {
			t.Errorf("Result[%d]: expected Name == %q, but got %q", i, exp.Name, l.Name)
		}
	}
}

func TestNewLeagueHandler(t *testing.T) {
	cdb, err := setUpCompetitionDB()
	if err != nil {
		t.Fatalf("unexpected error setting up competition db: %s", err.Error())
	}
	defer tearDownCompetitionDB(cdb)
	e := echo.New()
	form := make(url.Values)
	expectedName := "English Premier League"
	form.Set("name", expectedName)
	req := httptest.NewRequest(echo.POST, "/league/new", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := makeNewLeagueHandler(e, cdb)
	err = h(c)
	if err != nil {
		t.Fatalf("unexpected error in handler: %s", err.Error())
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected rec.Code == http.StatusOK, but got %d", rec.Code)
	}

	leagues, err := cdb.GetAllLeagues()
	if err != nil {
		t.Fatalf("unexpected error in GetAllLeauges: %s", err.Error())
	}
	lCount := len(leagues)
	if lCount != 1 {
		t.Fatalf("expected 1 league to exist, but found %d", lCount)
	}

	l := leagues[0]

	if l.Name != expectedName {
		t.Errorf("expected name to be %q, but got %q", expectedName, l.Name)
	}
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
