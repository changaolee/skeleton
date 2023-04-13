// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package middleware

import "github.com/gin-gonic/gin"

// AuthStrategy 身份认证策略，定义了身份认证的方法.
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator 用于切换不同的身份认证策略.
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy 用于设置身份认证策略.
func (o *AuthOperator) SetStrategy(strategy AuthStrategy) {
	o.strategy = strategy
}

// AuthFunc 执行身份认证.
func (o *AuthOperator) AuthFunc() gin.HandlerFunc {
	return o.strategy.AuthFunc()
}
