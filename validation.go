package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var (
	reservedNames = []string{"mule", "guest", "login", "logout"}
)

func IsValidName(testName string) bool {
	y := strings.ToLower(testName)
	for _, x := range reservedNames {
		if x == y {
			return false
		}
	}
	l := len(testName)
	if l < 3 {
		return false
	}
	for _, rn := range testName {
		if !unicode.In(rn, unicode.L, unicode.N) {
			return false
		}
	}
	return true
}

func IsValidPass(testpwd string) bool {
	l := len(testpwd)
	if l < 5 {
		return false
	}
	for _, rn := range testpwd {
		if !unicode.In(rn, unicode.L, unicode.N) {
			return false
		}
	}
	return true
}

func PassCheck(name, pwd string) bool {
	lvl := Level(name)
	if lvl == "" {
		return false
	}
	fileName := fmt.Sprintf(PASSFL, lvl, name)
	pwdfile, err := os.Open(fileName)
	if err != nil {
		ErrLog(err)
		return false
	}
	defer pwdfile.Close()
	pwdfileText := bufio.NewScanner(pwdfile)
	pwdfileText.Split(bufio.ScanWords)
	pwdfileText.Scan()
	if err = pwdfileText.Err(); err != nil {
		ErrLog(err)
		return false
	}
	testpwd := pwdfileText.Text()
	if !IsValidPass(testpwd) {
		err = ErrF("User %s has invalidly formatted password %q", name, testpwd)
		ErrLog(err)
		return false
	}
	return pwd == testpwd
}
