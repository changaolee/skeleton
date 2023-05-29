// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package apiserver

import (
	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/gin-gonic/gin"

	_ "github.com/changaolee/skeleton/internal/pkg/validator"

	"github.com/changaolee/skeleton/internal/apiserver/controller/v1/user"
	"github.com/changaolee/skeleton/internal/apiserver/store/mysql"
	"github.com/changaolee/skeleton/internal/pkg/middleware/auth"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) {
	jwtStrategy, _ := newJWTAuth().(auth.JWTStrategy)

	// 认证相关接口
	g.POST("/login", jwtStrategy.LoginHandler)     // 用户登录
	g.POST("/logout", jwtStrategy.LogoutHandler)   // 用户登出
	g.POST("/refresh", jwtStrategy.RefreshHandler) // 刷新 Token

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	// v1 路由分组
	storeIns, _ := mysql.GetMySQLInstance(nil)
	v1 := g.Group("/v1")
	{
		// 用户相关接口
		userv1 := v1.Group("/users")
		{
			userController := user.NewUserController(storeIns)

			userv1.POST("", userController.Create) // 创建用户

			// 权限检查中间件
			userv1.Use(auto.AuthFunc())

			// todo: 用户管理接口
			// userv1.PUT(":name", userController.Update)
			// userv1.GET("", userController.List)
			userv1.GET(":name", userController.Get)
		}
	}
}
