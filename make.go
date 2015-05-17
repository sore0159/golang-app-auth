package auth

import (
	"fmt"
	"io/ioutil"
	"os"
)

const GUESTLIMIT = 50

func NewUser(name string, pass string) *Data {
	lvl := Level(name)
	if lvl != "" {
		ErrLog("Tried to create existing user", name)
		return nil
	}
	u := MakeData(name, UNKNOWN)
	_, err := u.SaveDir("")
	//err := os.MkdirAll(new_dir, 0700)
	if err != nil {
		ErrLog("Failed to create user dir for", name, err)
		return nil
	}
	err = ioutil.WriteFile(u.PassFl(), []byte(pass), 0600)
	if err != nil {
		ErrLog("Failed to create pass file for", name, err)
		return nil
	}
	err = u.setupLoggers()
	if err != nil {
		ErrLog("Failed to set up Loggers for", name, err)
		return nil
	}
	return u
}

func NewGuest() *Data {
	d := MakeData(GUEST, GUEST)
	err := d.setupLoggers()
	if err != nil {
		ErrLog("GUEST CREATION FAILED: Failed to set up Guest Loggers", err)
		return nil
	}
	return d
}

func (d *Data) CreateGuestSaveDir() error {
	var i int
	var guestID string
	guestDIR := fmt.Sprintf(USERSDIR, "guest/")
	for i < GUESTLIMIT {
		guestID = fmt.Sprintf("guest-%03d", i)
		if _, err := os.Stat(guestDIR + guestID); os.IsNotExist(err) {
			break
		} else {
			i++
		}
		if i == GUESTLIMIT {
			return ErrF("Failed to create guest: limit reached!")
		}
	}
	d.UserName = guestID
	_, err := d.SaveDir("")
	//err := os.MkdirAll(new_dir, 0700)
	if err != nil {
		ErrLog(err)
		return err
	}
	err = d.Login(d.W, d.R)
	if err != nil {
		ErrLog(err)
		return err
	}
	return nil
}
