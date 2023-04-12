// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// base Error.
type base struct {
	msg string
	*stack
}

func (e *base) Error() string {
	return e.msg
}

func (e *base) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.msg)
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.msg)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.msg)
	}
}

// withCode Error.
type withCode struct {
	err   error
	code  int
	cause error
	*stack
}

func (w *withCode) Error() string { return fmt.Sprintf("%v", w) }

func (w *withCode) Cause() error { return w.cause }

func (w *withCode) Unwrap() error { return w.cause }

func (w *withCode) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})
		var jsonData []map[string]interface{}

		var (
			flagDetail bool
			flagTrace  bool
			modeJSON   bool
		)

		if state.Flag('#') {
			modeJSON = true
		}

		if state.Flag('-') {
			flagDetail = true
		}
		if state.Flag('+') {
			flagTrace = true
		}

		sep := ""
		errs := list(w)
		length := len(errs)
		for k, e := range errs {
			finfo := buildFormatInfo(e)
			jsonData, str = format(length-k-1, jsonData, str, finfo, sep, flagDetail, flagTrace, modeJSON)
			sep = "; "

			if !flagTrace {
				break
			}

			if !flagDetail && !flagTrace && !modeJSON {
				break
			}
		}
		if modeJSON {
			var byts []byte
			byts, _ = json.Marshal(jsonData)

			str.Write(byts)
		}

		_, _ = fmt.Fprintf(state, "%s", strings.Trim(str.String(), "\r\n\t"))
	default:
		finfo := buildFormatInfo(w)
		// Externally-safe error message
		_, _ = fmt.Fprintf(state, finfo.message)
	}
}

// withStack Error.
type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }

func (w *withStack) Unwrap() error {
	if e, ok := w.error.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
	}

	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// withMessage Error.
type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg }

func (w *withMessage) Cause() error { return w.cause }

func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
