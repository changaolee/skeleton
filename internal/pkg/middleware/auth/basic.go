// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package auth

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/internal/pkg/middleware"
	"github.com/changaolee/skeleton/pkg/errors"
)

type CompareFunc func(username string, password string) bool

// BasicStrategy 定义 Basic 认证策略.
type BasicStrategy struct {
	compare CompareFunc
}

var _ middleware.AuthStrategy = &BasicStrategy{}

// NewBasicStrategy 基于给定的 compare 方法创建一个 Basic 认证策略.
func NewBasicStrategy(compare CompareFunc) BasicStrategy {
	return BasicStrategy{compare}
}

// AuthFunc 定义 Basic 认证策略作为 Gin 中间件.
func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ["Basic", "xxx"]
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong"),
				nil,
			)
			c.Abort()

			return
		}

		// username:password
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		// [username, password]
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !b.compare(pair[0], pair[1]) {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong"),
				nil,
			)
			c.Abort()

			return
		}

		c.Set(middleware.UsernameKey, pair[0])

		c.Next()
	}
}
