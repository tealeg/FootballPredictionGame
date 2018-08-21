package app

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"

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

func makeCreateUserHandler(adb *user.AccountDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		isAdmin := strings.ToLower(c.FormValue("isadmin")) == "yes"
		acc := user.Account{
			Forename:       c.FormValue("forename"),
			Surname:        c.FormValue("surname"),
			Name:           c.FormValue("username"),
			HashedPassword: HashPassword(c.FormValue("password")),
			IsAdmin:        isAdmin,
		}
		err := adb.Create(acc)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.Redirect(http.StatusSeeOther, "/login")
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
	e.GET("/user/admin/exists", makeWelcomeHandler(adb))
	// e.GET("/firstuser", firstUserHandler)
	e.GET("/newuser", newUserHandler)
	e.POST("/createuser", makeCreateUserHandler(adb))
	e.GET("/login", loginHandler)
	e.POST("/authenticate", makeAuthenticationHandler(e, adb))
}
