package main

import (
	"log"

	"github.com/tealeg/FPG2/app"
	"github.com/tealeg/FPG2/user"
)

func main() {
	adb, err := user.NewAccountDB("user")
	if err != nil {
		log.Fatal(err)
	}
	app.Serve(":9090", adb)
}