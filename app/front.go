package app

import (
	"net/http"

	"github.com/tealeg/FPG2/user"
)

func makeFrontPageHandler(adb *user.AccountDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>
  <body>
    <h1>Welcome</h1>
  </body>
</html>
`
		w.Write([]byte(page))

	}
}

func setupFrontPageHandler(adb *user.AccountDB) {
	http.HandleFunc("/frontpage", SecurePage(adb, makeFrontPageHandler(adb)))
}
