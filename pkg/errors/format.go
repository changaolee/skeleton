// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package errors

import (
	"bytes"
	"fmt"
)

type formatInfo struct {
	code    int
	message string
	err     string
	stack   *stack
}

func list(e error) []error {
	var ret []error

	if e != nil {
		if w, ok := e.(interface{ Unwrap() error }); ok {
			ret = append(ret, e)
			ret = append(ret, list(w.Unwrap())...)
		} else {
			ret = append(ret, e)
		}
	}
	return ret
}

func format(k int, jsonData []map[string]interface{}, str *bytes.Buffer, finfo *formatInfo,
	sep string, flagDetail, flagTrace, modeJSON bool,
) ([]map[string]interface{}, *bytes.Buffer) {
	// nolint: nestif
	if modeJSON {
		data := map[string]interface{}{}
		if flagDetail || flagTrace {
			data = map[string]interface{}{
				"message": finfo.message,
				"code":    finfo.code,
				"error":   finfo.err,
			}

			caller := fmt.Sprintf("#%d", k)
			if finfo.stack != nil {
				f := Frame((*finfo.stack)[0])
				caller = fmt.Sprintf("%s %s:%d (%s)",
					caller,
					f.file(),
					f.line(),
					f.name(),
				)
			}
			data["caller"] = caller
		} else {
			data["error"] = finfo.message
		}
		jsonData = append(jsonData, data)
	} else {
		if flagDetail || flagTrace {
			if finfo.stack != nil {
				f := Frame((*finfo.stack)[0])
				_, _ = fmt.Fprintf(str, "%s%s - #%d [%s:%d (%s)] (%d) %s",
					sep,
					finfo.err,
					k,
					f.file(),
					f.line(),
					f.name(),
					finfo.code,
					finfo.message,
				)
			} else {
				_, _ = fmt.Fprintf(str, "%s%s - #%d %s", sep, finfo.err, k, finfo.message)
			}
		} else {
			_, _ = fmt.Fprintf(str, finfo.message)
		}
	}

	return jsonData, str
}

func buildFormatInfo(e error) *formatInfo {
	var finfo *formatInfo

	switch err := e.(type) {
	case *base:
		finfo = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.msg,
			err:     err.msg,
			stack:   err.stack,
		}
	case *withStack:
		finfo = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.Error(),
			err:     err.Error(),
			stack:   err.stack,
		}
	case *withCode:
		coder, ok := codes[err.code]
		if !ok {
			coder = unknownCoder
		}

		extMsg := coder.String()
		if extMsg == "" {
			extMsg = err.err.Error()
		}

		finfo = &formatInfo{
			code:    coder.Code(),
			message: extMsg,
			err:     err.err.Error(),
			stack:   err.stack,
		}
	default:
		finfo = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.Error(),
			err:     err.Error(),
		}
	}

	return finfo
}
