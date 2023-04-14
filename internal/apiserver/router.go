// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package apiserver

import (
	_ "github.com/changaolee/skeleton/internal/pkg/validator"
	"github.com/gin-gonic/gin"

	"github.com/changaolee/skeleton/internal/apiserver/store/mysql"
	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/internal/pkg/middleware/auth"
	"github.com/changaolee/skeleton/pkg/errors"

	"github.com/changaolee/skeleton/internal/apiserver/controller/v1/user"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) {
	jwtStrategy, _ := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	// v1 路由分组
	storeIns, _ := mysql.GetMySQLInstance(nil)
	v1 := g.Group("/v1")
	{
		// users 路由分组
		userv1 := v1.Group("/users")
		{
			userController := user.New(storeIns)

			userv1.POST("", userController.Create) // 创建用户
			userv1.Use(auto.AuthFunc())

			// todo
			// userv1.PUT(":name", userController.Update)
			// userv1.GET("", userController.List)
			// userv1.GET(":name", userController.Get) // admin api
		}
	}
}
