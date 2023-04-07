// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/changaolee/skeleton/pkg/log"
)

// UsernameKey 在 Gin 上下文中定义代表秘钥所有者的 Key.
const UsernameKey = "username"

// Context 是一个 Gin 中间件，将公共 Key 注入到上下文中.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.KeyRequestID, c.GetString(XRequestIDKey))
		c.Set(log.KeyUsername, c.GetString(UsernameKey))
		c.Next()
	}
}
