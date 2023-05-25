// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

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
