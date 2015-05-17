package auth

import (
	"fmt"
	"mule/mylog"
)

var Log = mylog.Log
var ErrLog = mylog.Err
var ErrF = mylog.ErrF

func (u *Data) Log(v ...interface{}) {
	u.logger(v...)
}

func (u *Data) ErrLog(v ...interface{}) {
	u.eLogger(v...)
}

func (u *Data) setupLoggers() error {
	u.eLogger = mylog.MakeErr(fmt.Sprintf("ERROR %s: ", u.Name), u.LogFl())
	u.logger = mylog.Make(fmt.Sprintf("USER %s: ", u.Name), u.LogFl())
	if u.eLogger == nil || u.logger == nil {
		return ErrF("Failed to make loggers")
	}
	return nil
}
