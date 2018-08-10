package app

import (
	"crypto/sha1"
	"fmt"
	"log"
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

func newUserHandler(w http.ResponseWriter, r *http.Request) {
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
        <input type="hidden" name="isadmin" value="no" />
        <input type="submit" value="Create account"/>
      </fieldset>
    </form>
  </body>
</html>
`))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	_, failed := r.URL.Query()["failed"]

	page := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>`
	if failed {
		page += `<span class="error">Login Failed</span>`
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
	w.Write([]byte(page))

}

func makeWelcomeHandler(adb *user.AccountDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exists, err := adb.AdminUserExists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/firstuser", http.StatusSeeOther)
		return
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
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func makeAuthenticationHandler(adb *user.AccountDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		account, err := adb.Get(r.FormValue("username"))
		if err != nil {
			log.Println("Couldn't get account: " + err.Error())
			w.Write([]byte(err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := HashPassword(r.FormValue("password"))
		if p != account.HashedPassword {
			log.Println("bad password")
			http.Redirect(w, r, "/login?failed=true", http.StatusSeeOther)
			return
		}
		cookie, err := GetAccountCookie(adb, *account)
		if err != nil {
			log.Println("Couldn't get cookie: " + err.Error())
			w.Write([]byte(err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Set Cookie")
		http.SetCookie(w, cookie)

		target := r.URL.Query().Get("target")
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

		log.Println("Redirect to " + target)
		w.Write([]byte(page))
		return
	}
}

func HashPassword(password string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}

func setupUserHandlers(adb *user.AccountDB) {
	http.HandleFunc("/", makeWelcomeHandler(adb))
	http.HandleFunc("/firstuser", firstUserHandler)
	http.HandleFunc("/newuser", newUserHandler)
	http.HandleFunc("/createuser", makeCreateUserHandler(adb))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/authenticate", makeAuthenticationHandler(adb))

}
