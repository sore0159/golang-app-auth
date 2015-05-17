package auth

import (
	"net/http"
)

func validLogin(r *http.Request) (*Data, bool) {
	formUser := r.FormValue("username")
	formPass := r.FormValue("password")
	if !PassCheck(formUser, formPass) {
		Log(formUser + " login failed with bad password")
		return nil, false
	}
	data := GetData(formUser)
	if data == nil {
		err := ErrF("Login attempt for %s, invalid user", formUser)
		Log(err)
		return nil, false
	}
	return data, true
}

func SignUp(r *http.Request) (data *Data, errNum int) {
	pass := r.FormValue("password")
	name := r.FormValue("username")
	if !IsValidName(name) || Level(name) != "" {
		errNum++
	}
	if !IsValidPass(pass) {
		errNum += 2
	}
	if errNum == 0 {
		data = NewUser(name, pass)
		if data == nil {
			ErrLog(ErrF("Sign Up failed even with valid name/pass:", name))
			errNum = 4
		}
	}
	return
}

func (d *Data) GoHome() {
	http.Redirect(d.W, d.R, d.HomeURL, http.StatusFound)
}

func (d *Data) GoMain() {
	http.Redirect(d.W, d.R, "/", http.StatusFound)
}
func (d *Data) GoGame() {
	http.Redirect(d.W, d.R, d.GameURL, http.StatusFound)
}

func RedirHome(d *Data) {
	http.Redirect(d.W, d.R, d.HomeURL, http.StatusFound)
}

func RedirMain(d *Data) {
	http.Redirect(d.W, d.R, "/", http.StatusFound)
}

func RedirGame(d *Data) {
	http.Redirect(d.W, d.R, d.GameURL, http.StatusFound)
}
