// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package core

import (
	"net/http"

	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/changaolee/skeleton/pkg/log"
	"github.com/gin-gonic/gin"
)

// ErrResponse 定义了发生错误时的返回消息.
type ErrResponse struct {
	Code      int    `json:"code"`                // 业务错误码
	Message   string `json:"message"`             // 对外展示的错误信息
	Reference string `json:"reference,omitempty"` // 解决此错误的参考文档
}

// WriteResponse 将错误或响应数据写入 HTTP 响应主体
// WriteResponse 使用 errors.ParseCoder 方法，根据错误类型，尝试从 err 中提取业务错误码和错误信息.
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
