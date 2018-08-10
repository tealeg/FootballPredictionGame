package app

import (
	"log"
	"net/http"

	"github.com/tealeg/FPG2/user"
)

func Serve(port string, adb *user.AccountDB) {
	setupUserHandlers(adb)
	setupFrontPageHandler(adb)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalf("Unexpected error in http.ListenAndServe: %s", err.Error())
	}

}
