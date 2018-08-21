package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo"
)

func TestMakeNewLeagueHandler(t *testing.T) {
	// adb := setupAccountDB()
	// defer tearDownAccountDB()
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
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected rec.Code == http.StatusSeeOther, but got %d", rec.Code)
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

	expectedLoc := fmt.Sprintf("/league/%d", l.ID)
	loc := rec.HeaderMap.Get("Location")
	if loc != expectedLoc {
		t.Errorf("expected location to be %q, but got %q", expectedLoc, loc)
	}

	if l.Name != expectedName {
		t.Errorf("expected name to be %q, but got %q", expectedName, l.Name)
	}

}
