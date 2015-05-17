package auth

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"time"
)

var secure = securecookie.New(hashKey, blockKey)

func (u *Data) Login(w http.ResponseWriter, r *http.Request) error {
	u.W = w
	u.R = r
	name := u.UserName
	var exp time.Time
	if u.Level == GUEST && len(name) > 7 {
		// AT THE STROKE OF MIDNIGHT
		name += time.Now().Round(time.Hour * 24).String()
		exp = time.Now().Round(time.Hour * 24).Add(time.Hour * 24)
	}
	secureValue, err := secure.Encode("username", name)
	if err != nil {
		u.ErrLog("Error logging in:", err)
		return err
	}
	cookie := &http.Cookie{
		Name:    "username",
		Value:   secureValue,
		Path:    "/",
		Expires: exp,
	}
	http.SetCookie(u.W, cookie)
	u.Log(u.UserName + "Logged in")
	return nil
}

func (u *Data) Logout() {
	if u.UserName == GUEST {
		return
	}
	delete(dataCache, u.UserName)
	GuestLogin(u.W)
	u.Log("Logged out")
}

func GuestLogin(w http.ResponseWriter) {
	Log("Logged in as Guest")
	cookie := &http.Cookie{
		Name:  "username",
		Value: GUEST,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func CookieUserName(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("username")
	if err == http.ErrNoCookie {
		GuestLogin(w)
		return GUEST
	} else if err != nil {
		ErrLog(err)
		GuestLogin(w)
		return GUEST
	}
	if cookie.Value == GUEST {
		return GUEST
	}
	var userName string
	if err = secure.Decode("username", cookie.Value, &userName); err != nil {
		ErrLog(err)
		GuestLogin(w)
		return GUEST
	}
	if len(userName) > 5 && userName[:6] == "guest-" {
		if userName == userName[:9]+time.Now().Round(time.Hour*24).String() {
			return userName[:9]
		} else {
			Log("Guest", userName[:9], "tried to log in with old cookie")
			GuestLogin(w)
			return GUEST
		}
	}
	return userName
}
