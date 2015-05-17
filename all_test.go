package auth

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	fmt.Println("TESTING")
}

var good_passwd_tests = []string{"aoeu0934", "99oaAOEU", "UAEUOUOO"}
var bad_passwd_tests = []string{"1234", "aou_UOou", "!!!!!"}

func TestPasswords(t *testing.T) {
	for _, test := range good_passwd_tests {
		if !IsValidPass(test) {
			err := ErrF("test: %s| want: pass| got: fail", test)
			t.Error(err)
		}
	}
	for _, test := range bad_passwd_tests {
		if IsValidPass(test) {
			err := ErrF("test: %s| want: fail| got: pass", test)
			t.Error(err)
		}
	}
}
