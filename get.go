package auth

import (
	"fmt"
	"net/http"
	"os"
)

var (
	dataCache  = make(map[string]*Data)
	USERLEVELS = []string{"admin", "known", "unknown", "guest"}
)

const APPNAME = "auth"

func CookieData(w http.ResponseWriter, r *http.Request) *Data {
	userName := CookieUserName(w, r)
	d := GetData(userName)
	if d == nil {
		GuestLogin(w)
		d = NewGuest()
		if d == nil {
			return nil
		}
	}
	u := d.Copy()
	u.W = w
	u.R = r
	return u
}

func GetData(userName string) *Data {
	var d *Data
	if d, _ = dataCache[userName]; d == nil {
		if userName == GUEST {
			d = NewGuest()
			if d == nil {
				return nil
			}
		} else {
			if d = Load(userName); d == nil {
				return nil
			}
		}
		d.Log("Loading into cache:", userName)
		dataCache[userName] = d
	}
	return d
}

func Load(userName string) *Data {
	lvl := Level(userName)
	if lvl == "" {
		ErrLog("Failed to find user", userName)
		return nil
	}
	d := MakeData(userName, lvl)
	err := d.setupLoggers()
	if err != nil {
		ErrLog("Failed to set up Loggers for", userName, err)
		return nil
	}
	return d
}

func Level(name string) string {
	if name == GUEST {
		return GUEST
	}
	for _, lvl := range USERLEVELS {
		if _, err := os.Stat(fmt.Sprintf(USERSDIR, lvl+"/"+name)); os.IsNotExist(err) {
			continue
		} else if err != nil {
			ErrLog(err)
			continue
		} else {
			return lvl
		}
	}
	return ""
}
