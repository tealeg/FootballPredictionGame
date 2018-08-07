package app

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"

	"github.com/tealeg/FPG2/user"
)

func firstUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>

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
        <input type="hidden" name="isadmin" value="yes" />
        <input type="submit" value="add user"/>
      </fieldset>
    </form>
  </body>
</html>
`))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>
    <form action="/checklogin" method="POST">
      <fieldset>
        <legend>Please login</legend>
        Username: <input type="text" name="username"/><br />
        Password: <input type="password" name="password"/><br />
      </fieldset>
    </form>
  </body>
</html>
`))

}

func makeWelcomeHandler(adb *user.AccountDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exists, err := adb.AdminUserExists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if exists {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		http.Redirect(w, r, "/firstuser", http.StatusSeeOther)
	}
}

func makeCreateUserHandler(adb *user.AccountDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		isAdmin := strings.ToLower(r.FormValue("isadmin")) == "yes"
		acc := user.Account{
			Forename:       r.FormValue("forename"),
			Surname:        r.FormValue("surname"),
			Name:           r.FormValue("username"),
			HashedPassword: HashPassword(r.FormValue("password")),
			IsAdmin:        isAdmin,
		}
		err := adb.Create(acc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func HashPassword(password string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}

func setupUserHandlers(adb *user.AccountDB) {
	http.HandleFunc("/", makeWelcomeHandler(adb))
	http.HandleFunc("/firstuser", firstUserHandler)
	http.HandleFunc("/createuser", makeCreateUserHandler(adb))
	http.HandleFunc("/login", loginHandler)

}
