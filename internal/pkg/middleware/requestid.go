// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// XRequestIDKey 用来定义 Gin 上下文中的键，代表请求的 uuid.
const XRequestIDKey = "X-Request-ID"

// RequestID 是一个 Gin 中间件，用来在每一个 HTTP 请求的 context, response 中注入 `X-Request-ID` 键值对.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否有 `X-Request-ID`，如果有则复用，没有则新建
		rid := c.GetHeader(XRequestIDKey)

		if rid == "" {
			rid = uuid.New().String()
			c.Request.Header.Set(XRequestIDKey, rid)
		}

		// 将 RequestID 保存在 gin.Context 中，方便后边程序使用
		c.Set(XRequestIDKey, rid)

		// 将 RequestID 保存在 HTTP 返回头中，Header 的键为 `X-Request-ID`
		c.Writer.Header().Set(XRequestIDKey, rid)
		c.Next()
	}
}
