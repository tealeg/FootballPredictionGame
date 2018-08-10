package app

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/tealeg/FPG2/user"
)

const UserCookieName = "FPG2UserName"

//
func GetAccountCookie(adb *user.AccountDB, acc user.Account) (*http.Cookie, error) {
	cookie := new(http.Cookie)
	cookie.Name = UserCookieName
	cookie.Value = acc.Name
	expiration := time.Now().Add(20 * time.Minute)
	acc.SessionExpires = expiration
	cookie.Expires = expiration

	err := adb.Update(acc.Name, acc)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func checkAccountCookie(e *echo.Echo, adb *user.AccountDB, c echo.Context, checkTime time.Time) bool {
	cookie, err := c.Cookie(UserCookieName)
	if err != nil {
		e.Logger.Error(err.Error())
		return false
	}
	e.Logger.Infof("Got cookie: %+v", cookie)
	acc, err := adb.Get(cookie.Value)
	if err != nil {
		e.Logger.Error(err.Error())
		return false
	}
	e.Logger.Infof("Got account: %+v", acc)
	// if acc.SessionExpires.Format(time.UnixDate) != c.RawExpires {
	// 	log.Println("date mismatch")
	// 	return false
	// }
	if acc.SessionExpires.UnixNano() < checkTime.UnixNano() {
		e.Logger.Warn("session expired")
		return false
	}
	return true
}

func SecurePage(e *echo.Echo, adb *user.AccountDB, h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !checkAccountCookie(e, adb, c, time.Now()) {
			e.Logger.Warn("Cookie check failed")
			return c.Redirect(http.StatusSeeOther, "/login?failed=timeout")
		}
		e.Logger.Info("Cookie check succeeded")
		return h(c)
	}
}
