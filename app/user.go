package app

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tealeg/FootballPredictionGame/user"
)

func newUserHandler(c echo.Context) error {
	return c.HTML(http.StatusOK, `<!DOCTYPE html>

<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>
    <form action="/createuser" method="POST">
      <fieldset>
        <legend>No admin user yet exists, please enter the admin users details here.</legend>
        Forename: <input type="text" name="forename"/><br />
        Surname: <input type="text" name="surname"/><br />
        Username: <input type="text" name="username"/><br />
        Password: <input type="password" name="password"/><br />
        <input type="hidden" name="isadmin" value="no" />
        <input type="submit" value="Create account"/>
      </fieldset>
    </form>
  </body>
</html>
`)
}

func loginHandler(c echo.Context) error {
	failed := c.QueryParam("failed")

	page := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>`
	switch failed {
	case "true":
		page += `<span class="error">Login Failed</span>`
	case "timeout":
		page += `<span class="error">Session expired</span>`
	}

	page += `<form action="/authenticate" method="POST">
      <fieldset>
        <legend>Please login</legend>
        Username: <input type="text" name="username"/><br />
        Password: <input type="password" name="password"/><br />
        <input type="submit" value="login"
      </fieldset>
    </form>
    <p>Not already a user? <a href="/newuser">Create an account.</a></p>
  </body>
</html>
`
	return c.HTML(http.StatusOK, page)

}

func makeAdminUserExistsHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		exists, err := adb.AdminUserExists()
		if err != nil {
			e.Logger.Error(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if exists {
			e.Logger.Error("Admin user exists")
			return c.JSON(http.StatusOK, true)
		}
		e.Logger.Error("Admin user does not exist")
		return c.JSON(http.StatusOK, false)
	}
}

type simpleResponse struct {
	Errors []string
}

func NewSimpleResponse() *simpleResponse {
	return &simpleResponse{Errors: []string{}}
}

func (s *simpleResponse) AddError(err error) {
	s.Errors = append(s.Errors, err.Error())
}

type createUserRequest struct {
	Forename string `json:"forename"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isadmin"`
}

func (cur *createUserRequest) CreateAccount(adb *user.AccountDB) error {
	acc := user.Account{
		Forename:       cur.Forename,
		Surname:        cur.Surname,
		Name:           cur.Username,
		HashedPassword: HashPassword(cur.Password),
		IsAdmin:        cur.IsAdmin,
	}
	return adb.Create(acc)
}

func (cur *createUserRequest) Validate(r *simpleResponse) error {
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

func makeCreateUserHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cur := new(createUserRequest)
		if err := c.Bind(cur); err != nil {
			e.Logger.Error(err.Error())
			return err
		}
		r := NewSimpleResponse()
		err := cur.Validate(r)
		if err != nil {
			for _, rerr := range r.Errors {
				e.Logger.Error(rerr)
			}
			return echo.NewHTTPError(http.StatusBadRequest, r.Errors)
		}
		err = cur.CreateAccount(adb)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, r)
	}
}

type loginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

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

func makeAuthenticationHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	e.Logger.Error("Creating AuthenticationHandler")
	return func(c echo.Context) error {
		e.Logger.Error("Authenticating")
		lr := new(loginRequest)
		if err := c.Bind(lr); err != nil {
			e.Logger.Error(err.Error())
			return err
		}
		e.Logger.Error(fmt.Sprintf("%+v", lr))
		r := NewSimpleResponse()
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
			e.Logger.Error("bad password")
			r.AddError(errors.New("Bad credentials - user name and password not valid for this service"))
			return c.JSON(http.StatusUnauthorized, *r)
		}
		cookie, err := GetAccountCookie(adb, *account)
		if err != nil {
			e.Logger.Error("Couldn't get cookie: " + err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		e.Logger.Error("Set Cookie")
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, *r)
	}
}

func makeLogOutHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := ExpireAccountCookie(e, c, adb)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		e.Logger.Error("expire cookie")
		c.SetCookie(cookie)
		return c.JSON(http.StatusOK, nil)
	}
}

func HashPassword(password string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}

func setupUserHandlers(e *echo.Echo, adb *user.AccountDB) {
	e.PUT("/authenticate", makeAuthenticationHandler(e, adb))
	e.GET("/user/admin/exists.json", makeAdminUserExistsHandler(e, adb))
	// e.GET("/newuser", newUserHandler)
	e.PUT("/user/new.json", makeCreateUserHandler(e, adb))
	// e.GET("/login", loginHandler)

	e.POST("/logout", makeLogOutHandler(e, adb))
}
