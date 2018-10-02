package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
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
func TestAdminUserExistsNoAdmin(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/user/admin/exists.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeAdminUserExistsHandler(e, adb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotZero(t, rec.Body)
	var result bool
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.False(t, result)
}

// When an admin user exists the AdminUserExistsHandler indicates this.
func TestAdminUserExistsWithAdmin(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)
	acc := user.Account{
		Forename:       "Geoff",
		Surname:        "Teale",
		Name:           "tealeg",
		HashedPassword: "IAmHashed",
		IsAdmin:        true,
	}
	err = adb.Create(acc)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/user/admin/exists.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeAdminUserExistsHandler(e, adb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotZero(t, rec.Body)

	var result bool
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.True(t, result)
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
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		aLen := len(r.Errors)
		if aLen != eLen {
			t.Fatalf("Case %d: expected %d errors, but got %d", i, eLen, aLen)
		}

		for j, msg := range r.Errors {
			expected := exp.Expected[j]
			assert.Equal(t, expected, msg)
		}
	}
}

// Accounts are created and match their requests.
func TestCreateAccount(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)

	cur := &createAccountRequest{
		Forename: "Bob",
		Surname:  "Bobfrey",
		Username: "bobit",
		Password: "lorena",
		IsAdmin:  true,
	}

	err = cur.createAccount(adb)
	assert.NoError(t, err)
	acc, err := adb.Get("bobit")
	assert.NoError(t, err)

	assert.Equal(t, acc.Forename, cur.Forename)
	assert.Equal(t, acc.Surname, cur.Surname)
	assert.Equal(t, acc.Name, cur.Username)
	assert.Equal(t, acc.HashedPassword, HashPassword(cur.Password))
	assert.Equal(t, acc.IsAdmin, cur.IsAdmin)
}

// The returned createAccountHandler creates accounts
func TestCreateAccountHandler(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user/new.json",
		strings.NewReader(`{"Forename":"bob","Surname":"bobfrey","Username":"bobit","Password":"lorena","IsAdmin":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeCreateAccountHandler(e, adb)
	err = h(c)
	assert.NoError(t, err)
	// TODO - maybe test that we actually create something?
}

// The returned createAccountHandler returns a list of validation errors
func TestCreateAccountHandlerWithBadRequest(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user/new.json",
		strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeCreateAccountHandler(e, adb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	errors := []string{}
	err = json.Unmarshal(rec.Body.Bytes(), &errors)
	assert.NoError(t, err)
	assert.Len(t, errors, 4)
}

// loginRequest.Validate indicates all validation failures.
func TestLoginRequestValidate(t *testing.T) {

	cases := []struct {
		LR          *loginRequest
		Expectation []string
	}{
		{
			LR: &loginRequest{},
			Expectation: []string{
				"No user name was provided",
				"No password was provided",
			},
		},
	}

	adb, err := setUpAccountDB()
	if err != nil {
		t.Fatalf("unexpected error creating AccountDB: %s", err.Error())
	}
	defer tearDownAccountDB(adb)

	for _, c := range cases {
		eLen := len(c.Expectation)
		sr := &simpleResponse{}
		err := c.LR.Validate(sr)
		if eLen > 0 {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		for j, e := range c.Expectation {
			er := sr.Errors[j]
			assert.Equal(t, e, er)

		}
	}
}

// With valid credentials, we can be authenticated
func TestAuthenticationHandlerWithValidCredetentials(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)

	acc := user.Account{
		Forename:       "bob",
		Surname:        "bobington",
		Name:           "bobit",
		HashedPassword: HashPassword("lorena"),
		IsAdmin:        false,
	}
	err = adb.Create(acc)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/authenticate",
		strings.NewReader(`{"Username":"bobit", "Password":"lorena"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeAuthenticationHandler(e, adb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotZero(t, rec.Body)

	cookies, ok := rec.HeaderMap["Set-Cookie"]
	assert.True(t, ok)

	cookie := cookies[0]
	fields := strings.Split(cookie, ";")
	parts := strings.Split(fields[0], "=")
	assert.Equal(t, "FPG2UserName", parts[0])
	assert.Equal(t, "bobit", parts[1])
}

func TestIsAdminUserFalse(t *testing.T) {
	adb, err := setUpAccountDB()
	assert.NoError(t, err)
	defer tearDownAccountDB(adb)
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/user/isadmin.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := makeIsAdminHandler(e, adb)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotZero(t, rec.Body)

	var result bool
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.False(t, result)
}
