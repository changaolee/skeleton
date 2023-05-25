package util

import (
	"fmt"
	"os"
	"strings"
)

const (
	DefaultErrorExitCode = 1
)

// fatal 打印 msg 并退出.
func fatal(msg string, code int) {
	if len(msg) > 0 {
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		_, _ = fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

var fatalErrHandler = fatal

func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

func checkErr(err error, handleErr func(string, int)) {
	if err == nil {
		return
	}

	msg := err.Error()
	if !strings.HasPrefix(msg, "error: ") {
		msg = fmt.Sprintf("error: %s", msg)
	}
	handleErr(msg, DefaultErrorExitCode)
}
