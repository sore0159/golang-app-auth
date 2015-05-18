package auth

import (
	//"fmt"
	//"html/template"
	"mule/helpers"
)

//const TEMP_TEMP = "templates/auth/%s.html"

var (
	mixTem  = helpers.GenTMixer(APPNAME)
	loginT  = mixTem("frame", "login")
	SignUpT = mixTem("frame", "signup")
)

/*
func mixTem(tmpls ...string) *template.Template {
	new_list := make([]string, 0)
	for _, val := range tmpls {
		new_list = append(new_list, fmt.Sprintf(TEMP_TEMP, val))
	}
	t, err := template.New("t").ParseFiles(new_list...)
	if err == nil {
		return t
	}
	panic(err)
}
*/

func LoginPage(d *Data) {
	if d.R.Method == "POST" {
		if data, ok := validLogin(d.R); ok {
			err := data.Login(d.W, d.R)
			if err == nil {
				d.GoHome()
				return
			}
		}
		// ERROR FLASH MESSAGE
	}
	formUser := d.R.FormValue("username")
	// maybe use session cookie for "last attempted login"
	d.ExeT(loginT, "frame", formUser)
}

func SignUpPage(d *Data) {
	pageData := make(map[string]string)
	if d.R.Method == "POST" {
		data, errNum := SignUp(d.R)
		switch errNum {
		case 0:
			err := data.Login(d.W, d.R)
			if err == nil {
				data.GoHome()
				return
			} else {
				ErrLog("Login Error for newly created account:", err)
				d.GoMain()
				return
			}
		case 1:
			pageData["password"] = d.R.FormValue("password")
			pageData["errors"] = "Invalid username!"
		case 2:
			pageData["username"] = d.R.FormValue("username")
			pageData["errors"] = "Invalid password!"
		default:
			pageData["errors"] = "Invalid everything!"
			// bad both!
		}
	}
	// maybe use session cookie for "last attempted login"
	d.ExeT(SignUpT, "frame", pageData)
}

func LogoutPage(d *Data) {
	GuestLogin(d.W)
	d.GoMain()
}

func redirHome(d *Data) {
	d.GoHome()
}
func redirGame(d *Data) {
	d.GoGame()
}
