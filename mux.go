package auth

import (
	"github.com/gorilla/mux"
	"net/http"
)

func SetupMux(r *mux.Router) {
	r.HandleFunc("/register", wrap(SignUpPage))
	r.HandleFunc("/login", wrap(LoginPage))
	r.HandleFunc("/logout", DataWrap(LogoutPage))
}

func wrap(f func(*Data)) func(http.ResponseWriter, *http.Request) {
	return DataWrap(BoolWrap(IsAuth, f))
}

func DataWrap(f func(*Data)) func(http.ResponseWriter, *http.Request) {
	g := func(w http.ResponseWriter, r *http.Request) {
		data := CookieData(w, r)
		if data == nil {
			err := ErrF("Can't load page %s, can't get userdata!", r.URL.Path)
			ErrLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.SetGame("")
		f(data)
	}
	return g
}

func BoolWrap(b func(*Data) bool, f func(*Data)) func(*Data) {
	g := func(d *Data) {
		if !b(d) {
			d.Log("Rejected from", d.R.URL.Path)
			d.GoHome()
			return
		}
		f(d)
	}
	return g
}

func IsAuth(d *Data) bool {
	if d.Level != GUEST {
		//"The menagarie is for guests only."
		return false
	}
	return true
}

func GameURLWrap(appName string, f func(*Data)) func(*Data) {
	g := func(d *Data) {
		d.SetGame(appName)
		f(d)
	}
	return g
}
