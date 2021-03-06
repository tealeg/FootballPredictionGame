package app

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tealeg/FootballPredictionGame/user"
)

//makeIsAdminHandler returns a handler that will indicate if the
//currently logged in user is and administrator.
func makeIsAdminHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		acc, err := GetUserAccount(c, adb)
		if err != nil {
			if err.Error() == "http: named cookie not present" {
				return c.JSON(http.StatusOK, false)
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}
		return c.JSON(http.StatusOK, acc.IsAdmin)
	}

}

// makeAdminUserExistsHandler returns a handler that will indicate if
// an admin user has been created already.  This request should not
// require authentication.
func makeAdminUserExistsHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		e.Logger.Debugf("Check if admin user exists")
		exists, err := adb.AdminUserExists()
		if err != nil {
			e.Logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if exists {
			e.Logger.Debugf("Admin user exists")
			return c.JSON(http.StatusOK, true)
		}
		e.Logger.Debugf("Admin user does not exist")
		return c.JSON(http.StatusOK, false)
	}
}

// createAccountRequest is a holder for data passed into new user requests.
type createAccountRequest struct {
	Forename string `json:"forename"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isadmin"`
}

// createAccount creates a user.Account based on its createAccountRequest
func (cur *createAccountRequest) createAccount(adb *user.AccountDB) error {
	acc := user.Account{
		Forename:       cur.Forename,
		Surname:        cur.Surname,
		Name:           cur.Username,
		HashedPassword: HashPassword(cur.Password),
		IsAdmin:        cur.IsAdmin,
	}
	return adb.Create(acc)
}

// Validate checks the members of a createAccountRequest for validity
// and populates a simpleResponse with the errors it finds.  The last
// error found will be returned, and can be used to indicate overall
// validation failure (or, if nil, success).
func (cur *createAccountRequest) Validate(r *simpleResponse) error {
	var err error
	if cur.Forename == "" {
		err = errors.New("Forename is empty")
		r.AddError(err)
	}
	if cur.Surname == "" {
		err = errors.New("Surname is empty")
		r.AddError(err)
	}
	if cur.Username == "" {
		err = errors.New("Username is empty")
		r.AddError(err)
	}
	if cur.Password == "" {
		err = errors.New("Password is empty")
		r.AddError(err)
	}
	return err
}

// makeCreateAccountHandler returns ar handler that will attempt to
// create a new user.Account based on the details provided.
func makeCreateAccountHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		e.Logger.Info("Creating user")
		cur := new(createAccountRequest)
		if err := c.Bind(cur); err != nil {
			e.Logger.Error(err.Error())
			return err
		}
		r := newSimpleResponse()
		err := cur.Validate(r)
		if err != nil {
			for _, rerr := range r.Errors {
				e.Logger.Warn(rerr)
			}
			return c.JSON(http.StatusBadRequest, r.Errors)
		}
		err = cur.createAccount(adb)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		e.Logger.Info("User created: %+v", r)
		return c.JSON(http.StatusOK, r)
	}
}

// loginRequest is a holder for data passed by a login request
type loginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Validate will validate the contents of a loginRequest.  Any errors
// will be recorded on the provided simpleResponse, and the final
// error will be returned, this can be used to indicate overall
// failure (or, if null, success).
func (lr *loginRequest) Validate(r *simpleResponse) error {
	var err error
	if lr.UserName == "" {
		err = errors.New("No user name was provided")
		r.AddError(err)
	}
	if lr.Password == "" {
		err = errors.New("No password was provided")
		r.AddError(err)
	}
	return err
}

// makeAuthenticationHandler returns a handler that can be used to
// process and approve or reject authentication requests.  If the
// request is succesful a cookie will be set and then handlers wrapped
// with the app/cookie.SecurePage middleware will gate their usage by
// checking for the presence and actuality of this cookie.
func makeAuthenticationHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	e.Logger.Debugf("Creating AuthenticationHandler")
	return func(c echo.Context) error {
		e.Logger.Debugf("Authenticating")
		lr := new(loginRequest)
		if err := c.Bind(lr); err != nil {
			e.Logger.Error(err.Error())
			return err
		}
		r := newSimpleResponse()
		err := lr.Validate(r)
		if err != nil {
			e.Logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, *r)
		}
		account, err := adb.Get(lr.UserName)
		if err != nil {
			e.Logger.Error("Couldn't get account: " + err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		p := HashPassword(lr.Password)
		if p != account.HashedPassword {
			e.Logger.Infof("bad password")
			r.AddError(errors.New("Bad credentials - user name and password not valid for this service"))
			return c.JSON(http.StatusUnauthorized, *r)
		}
		cookie, err := GetAccountCookie(adb, *account)
		if err != nil {
			e.Logger.Error("Couldn't get cookie: " + err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		e.Logger.Debugf("Set Cookie")
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, *r)
	}
}

// makeLogOutHandler returns a handler that will expire a users
// session cookie, and thus require them to reauthenticate before
// accessing any handler wrapped with the app/cookie.SecurePage
// middleware.
func makeLogOutHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := ExpireAccountCookie(e, c, adb)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		e.Logger.Debugf("expire cookie")
		c.SetCookie(cookie)
		return c.JSON(http.StatusOK, nil)
	}
}

// HashPassword will return a one-way hash (sha1) of the provided password.
func HashPassword(password string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}

// setupUserHandlers binds handlers to all endpoints concerning user
// management and authentication
func setupUserHandlers(e *echo.Echo, adb *user.AccountDB) {
	e.PUT("/authenticate", makeAuthenticationHandler(e, adb))
	e.GET("/user/admin/exists.json", makeAdminUserExistsHandler(e, adb))
	e.GET("/user/isadmin.json", SecurePage(e, adb, makeIsAdminHandler(e, adb)))
	e.POST("/user/new.json", makeCreateAccountHandler(e, adb))
	e.GET("/logout", SecurePage(e, adb, makeLogOutHandler(e, adb)))
}
