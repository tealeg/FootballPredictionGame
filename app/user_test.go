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

// setUpAccountDB is a utility function to simplify creating a
// user.AccountDB for testing purposes.  It should be used in
// conjuction with tearDownAccountDB.
func setUpAccountDB() (*user.AccountDB, error) {
	return user.NewAccountDB("test-account")
}

// tearDownAccountDB is a utility function, to simplify clearing up a
// user.AccountDB used in a test. It should be used in conjuction with
// setUpAccountDB.
func tearDownAccountDB(adb *user.AccountDB) {
	adb.Close()
	os.Remove("test-account.db")
}

// When no admin user yet exists the AdminUserExistsHandler indicates this.
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

// When an admin user exists the AdminUserExistsHandler indicates this.
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

// createAccountRequests can be validated and will indicate which any
// and all fields that have issues, by means of a simpleResponse.
// Overall failues will be indicated by the return value of the
// Validate function.
func TestCreateAccountRequestValidation(t *testing.T) {
	var expectations = []struct {
		CAR      *createAccountRequest
		Expected []string
	}{
		{
			CAR: &createAccountRequest{},
			Expected: []string{
				"Forename is empty",
				"Surname is empty",
				"Username is empty",
				"Password is empty",
			},
		},
		{
			CAR: &createAccountRequest{
				Forename: "Bob",
			},
			Expected: []string{
				"Surname is empty",
				"Username is empty",
				"Password is empty",
			},
		},
		{
			CAR: &createAccountRequest{
				Surname: "Bobfrey",
			},
			Expected: []string{
				"Forename is empty",
				"Username is empty",
				"Password is empty",
			},
		},
		{
			CAR: &createAccountRequest{
				Username: "bobit",
			},
			Expected: []string{
				"Forename is empty",
				"Surname is empty",
				"Password is empty",
			},
		},
		{
			CAR: &createAccountRequest{
				Password: "lorena",
			},
			Expected: []string{
				"Forename is empty",
				"Surname is empty",
				"Username is empty",
			},
		},
		{
			CAR: &createAccountRequest{
				Forename: "Bob",
				Surname:  "Bobfrey",
				Username: "bobit",
				Password: "lorena",
			},
			Expected: []string{},
		},
	}

	for i, exp := range expectations {
		r := &simpleResponse{}
		eLen := len(exp.Expected)
		err := exp.CAR.Validate(r)
		if eLen > 0 {
			if err == nil {
				t.Fatalf("Case %d: expected error in validation, but got none", i)
			}
		} else {
			if err != nil {
				t.Fatalf("Case %d: unexpected error in validation: %s", i, err.Error())
			}
		}
		aLen := len(r.Errors)
		if aLen != eLen {
			t.Fatalf("Case %d: expected %d errors, but got %d", i, eLen, aLen)
		}

		for j, msg := range r.Errors {
			expected := exp.Expected[j]
			if msg != expected {
				t.Errorf("Case %d: r.Errors[%d] == %q, but should be %q", i, j, msg, expected)
			}
		}
	}
}
