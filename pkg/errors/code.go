// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var unknownCoder = defaultCoder{
	0,
	http.StatusInternalServerError,
	"An internal server error occurred",
	"https://github.com/changaolee/skeleton/tree/main/pkg/errors",
}

// Coder 定义了错误码详细信息的接口.
type Coder interface {
	Code() int         // 返回错误码
	HTTPStatus() int   // 返回与错误码相对应的 HTTP 状态码
	String() string    // 返回展示给外部用户的错误文本
	Reference() string // 返回给用户的详细文档
}

type defaultCoder struct {
	C    int    // 错误码
	HTTP int    // HTTP 状态码
	Ext  string // 错误文本
	Ref  string // 引用文档
}

func (coder defaultCoder) Code() int { return coder.C }

func (coder defaultCoder) String() string { return coder.Ext }

func (coder defaultCoder) HTTPStatus() int {
	if coder.HTTP == 0 {
		return http.StatusInternalServerError
	}
	return coder.HTTP
}

func (coder defaultCoder) Reference() string { return coder.Ref }

// codes 存储错误码到 Coder 的映射.
var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

// Register 注册一个用户定义的错误码（会覆盖已存在的错误码）.
func Register(coder Coder) {
	checkErrorCode(coder.Code())

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

// MustRegister 注册一个用户定义的错误码（与已存在错误码冲突时会 panic）.
func MustRegister(coder Coder) {
	checkErrorCode(coder.Code())

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

func checkErrorCode(code int) {
	if code == 0 {
		panic("code `0` is reserved by `github.com/changaolee/skeleton/errors` as unknownCode error code")
	}
}

// ParseCoder 将 err 映射为 withCode 错误并解析其中错误码.
// 解析失败会返回 unknown Error.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}
	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}
	return unknownCoder
}
