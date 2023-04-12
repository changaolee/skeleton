package code

import (
	"net/http"

	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/novalagung/gubrak"
)

type ErrCode struct {
	C    int    // 错误码
	HTTP int    // HTTP 状态码
	Ext  string // 错误文本
	Ref  string // 引用文档
}

var _ errors.Coder = &ErrCode{}

func (e ErrCode) Code() int {
	return e.C
}

func (e ErrCode) HTTPStatus() int {
	if e.HTTP == 0 {
		return http.StatusInternalServerError
	}
	return e.HTTP
}

func (e ErrCode) String() string {
	return e.Ext
}

func (e ErrCode) Reference() string {
	return e.Ref
}

func register(code int, httpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}

	errors.MustRegister(coder)
}
