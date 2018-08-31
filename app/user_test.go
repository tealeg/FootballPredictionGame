package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/tealeg/FootballPredictionGame/user"
)

func setUpAccountDB() (*user.AccountDB, error) {
	return user.NewAccountDB("test-account")
}

func tearDownAccountDB(adb *user.AccountDB) {
	adb.Close()
	os.Remove("test-account.db")
}

func TestAdminUserExists(t *testing.T) {
	adb, err := setUpAccountDB()
	if err != nil {
		t.Fatalf("Unexpected error in setUpAccountDB: %s", err.Error())
	}
	defer tearDownAccountDB(adb)
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/user/admin/exists.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeAdminUserExistsHandler(e, adb)
	err = h(c)
	if err != nil {
		t.Fatalf("Unexpected error in handler: %s", err.Error())
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected rec.Code = OK, but got %s", http.StatusText(rec.Code))
	}
	if rec.Body.Len() == 0 {
		t.Error("Empty body")
	}
	var result bool
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("error unmarshalling body: %s", err.Error())
	}
	if result {
		t.Error("expected handler to return false, but got true")
	}

}

func TestAdminUserExistsWithAdmin(t *testing.T) {
	adb, err := setUpAccountDB()
	if err != nil {
		t.Fatalf("Unexpected error in setUpAccountDB: %s", err.Error())
	}
	defer tearDownAccountDB(adb)
	acc := user.Account{
		Forename:       "Geoff",
		Surname:        "Teale",
		Name:           "tealeg",
		HashedPassword: "IAmHashed",
		IsAdmin:        true,
	}
	err = adb.Create(acc)
	if err != nil {
		t.Fatalf("unexpected error creating account: %s", err.Error())
	}

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/user/admin/exists.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeAdminUserExistsHandler(e, adb)
	err = h(c)
	if err != nil {
		t.Fatalf("Unexpected error in handler: %s", err.Error())
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected rec.Code = OK, but got %s", http.StatusText(rec.Code))
	}
	if rec.Body.Len() == 0 {
		t.Error("Empty body")
	}
	var result bool
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("error unmarshalling body: %s", err.Error())
	}
	if !result {
		t.Error("expected handler to return true, but got false")
	}

}
