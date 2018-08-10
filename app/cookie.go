package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tealeg/FPG2/user"
)

const UserCookieName = "FPG2UserName"

//
func GetAccountCookie(adb *user.AccountDB, acc user.Account) (*http.Cookie, error) {
	expiration := time.Now().Add(20 * time.Minute)
	acc.SessionExpires = expiration
	err := adb.Update(acc.Name, acc)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{Name: UserCookieName, Value: acc.Name, Expires: expiration, RawExpires: expiration.Format(time.UnixDate), Domain: "www.teale.de"}, nil
}

func checkAccountCookie(adb *user.AccountDB, r *http.Request, checkTime time.Time) bool {
	c, err := r.Cookie(UserCookieName)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	log.Println(fmt.Sprintf("Got cookie: %+v", c))
	acc, err := adb.Get(c.Value)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	log.Println(fmt.Sprintf("Got account: %+v", acc))
	if acc.SessionExpires.Format(time.UnixDate) != c.RawExpires {
		log.Println("date mismatch")
		return false
	}
	if acc.SessionExpires.UnixNano() < checkTime.UnixNano() {
		log.Println("session expired")
		return false
	}
	return true
}

func SecurePage(adb *user.AccountDB, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAccountCookie(adb, r, time.Now()) {
			log.Println("Cookie check failed")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		log.Println("Cookie check succeeded")
		h(w, r)
		return
	}
}
