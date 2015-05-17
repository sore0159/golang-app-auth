package auth

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Data struct {
	UserName string
	Name     string
	Level    string
	HomeURL  string
	GameURL  string
	logger   func(...interface{})
	eLogger  func(...interface{})
	W        http.ResponseWriter
	R        *http.Request
}

const (
	ADMIN    = "admin"
	KNOWN    = "known"
	UNKNOWN  = "unknown"
	GUEST    = "guest"
	GUESTLOG = "data/logs/guestlogs.txt"
	USERSDIR = "data/users/%s"
	HOMEURL  = "/%s"
	GAMEURL  = "/%s/%s"
	PASSFL   = USERSDIR + "/%s/password.txt"
	LOGFL    = USERSDIR + "/%s/logs.txt"
	SAVESDIR = USERSDIR + "/%s/saves/%s"
)

func MakeData(name, level string) *Data {
	d := &Data{
		UserName: name,
		Level:    level,
	}
	if level != GUEST {
		d.Name = name
	} else {
		d.Name = GUEST
	}
	d.HomeURL = fmt.Sprintf(HOMEURL, d.Name)
	return d
}

func (d *Data) SaveDir(gameName string) (string, error) {
	if d.Level == GUEST && d.UserName == GUEST {
		err := d.CreateGuestSaveDir()
		if err != nil {
			ErrLog("Couldn't make guest savedir", err)
			return "", err
		}
	}
	dirName := fmt.Sprintf(SAVESDIR, d.Level, d.UserName, gameName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0700)
		if err != nil {
			ErrLog("Couldn't make savedir", dirName, err)
			return "", err
		}
	} else if err != nil {
		ErrLog("Couldn't stat savedir", dirName, err)
		return "", err
	}
	return dirName, nil
}

func (d *Data) PassFl() string {
	return fmt.Sprintf(PASSFL, d.Level, d.UserName)
}

func (d *Data) LogFl() string {
	if d.Level == GUEST {
		return GUESTLOG
	} else {
		return fmt.Sprintf(LOGFL, d.Level, d.UserName)
	}
}

func (d *Data) Copy() *Data {
	return &Data{
		UserName: d.UserName,
		Name:     d.Name,
		Level:    d.Level,
		HomeURL:  d.HomeURL,
		GameURL:  d.GameURL,
		logger:   d.logger,
		eLogger:  d.eLogger,
	}
}

func (d *Data) SetGame(appName string) {
	d.GameURL = fmt.Sprintf(GAMEURL, d.Name, appName)
}

type pageDat struct {
	HomeURL string
	GameURL string
	Name    string
	Level   string
	AppDat  interface{}
}

func (d *Data) NewPageDat(data interface{}) *pageDat {
	return &pageDat{
		HomeURL: d.HomeURL,
		GameURL: d.GameURL,
		Name:    d.Name,
		Level:   d.Level,
		AppDat:  data,
	}
}
func (d *Data) ExeT(t *template.Template, name string, appData interface{}) {

	err := t.ExecuteTemplate(d.W, name, d.NewPageDat(appData))
	if err != nil {
		d.ErrLog(err)
	}
}
