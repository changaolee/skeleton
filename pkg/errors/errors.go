// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package errors

import "fmt"

// New 基于给定的 message 返回一个 base 错误.
func New(message string) error {
	return &base{
		msg:   message,
		stack: callers(),
	}
}

// Errorf 基于给定的参数格式化字符串返回一个 base 错误.
func Errorf(format string, args ...interface{}) error {
	return &base{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// WithCode 创建一个 withCode 错误.
func WithCode(code int, format string, args ...interface{}) error {
	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		stack: callers(),
	}
}

// Wrap 基于给定的 message 包装错误.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   fmt.Errorf(message),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

// Wrapf 基于给定的参数格式化字符串包装错误.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   fmt.Errorf(format, args...),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
	}
}

// WrapC 将 err 封装成一个 withCode 错误.
func WrapC(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		cause: err,
		stack: callers(),
	}
}

// Cause 返回错误的根本原因（如果 err 实现了 Cause() 方法）.
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(interface{ Cause() error })
		if !ok {
			break
		}
		if cause.Cause() == nil {
			break
		}
		err = cause.Cause()
	}
	return err
}

// IsCode 判断 err 链中是否有错误码为 code 的错误.
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}
		if v.cause != nil {
			return IsCode(v.cause, code)
		}
		return false
	}
	return false
}
