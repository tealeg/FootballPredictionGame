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

func makeWelcomeHandler(adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		exists, err := adb.AdminUserExists()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if exists {
			return c.JSON(http.StatusOK, true)
		}
		return c.JSON(http.StatusOK, false)
	}
}

type createUserResponse struct {
	Errors []string
}

func (cur *createUserResponse) AddError(err error) {
	cur.Errors = append(cur.Errors, err.Error())
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

func (cur *createUserRequest) Validate(r *createUserResponse) error {
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

func makeCreateUserHandler(adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cur := new(createUserRequest)
		if err := c.Bind(cur); err != nil {
			return err
		}
		r := createUserResponse{}
		err := cur.Validate(&r)
		if err != nil {
			return c.JSON(http.StatusOK, r)
		}
		err = cur.CreateAccount(adb)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, r)
	}
}

func makeAuthenticationHandler(e *echo.Echo, adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		account, err := adb.Get(c.FormValue("username"))
		if err != nil {
			e.Logger.Warn("Couldn't get account: " + err.Error())
			return c.Redirect(http.StatusSeeOther, "/login?failed=true")
		}

		p := HashPassword(c.FormValue("password"))
		if p != account.HashedPassword {
			e.Logger.Warn("bad password")
			return c.Redirect(http.StatusSeeOther, "/login?failed=true")
		}
		cookie, err := GetAccountCookie(adb, *account)
		if err != nil {
			e.Logger.Warn("Couldn't get cookie: " + err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		e.Logger.Info("Set Cookie")
		c.SetCookie(cookie)

		target := c.QueryParam("target")
		if target == "" {
			target = "/frontpage"
		}

		page := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		  <head>
		    <meta http-equiv="refresh" content="0; url=%s" />
		  </head>
		  <body>
		    <p>Redirecting, please wait...</p>
		  </body>
		</html>
		`, target)

		e.Logger.Info("Redirect to " + target)
		return c.HTML(http.StatusOK, page)
	}
}

func HashPassword(password string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}

func setupUserHandlers(e *echo.Echo, adb *user.AccountDB) {
	e.GET("/user/admin/exists.json", makeWelcomeHandler(adb))
	// e.GET("/newuser", newUserHandler)
	e.POST("/user/new.json", makeCreateUserHandler(adb))
	e.GET("/login", loginHandler)
	e.POST("/authenticate", makeAuthenticationHandler(e, adb))
}
