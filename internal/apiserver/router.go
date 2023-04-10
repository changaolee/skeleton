// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package apiserver

import (
	"github.com/changaolee/skeleton/internal/apiserver/controller/v1/user"
	"github.com/changaolee/skeleton/internal/apiserver/store/mysql"
	"github.com/gin-gonic/gin"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) {
	storeIns, _ := mysql.GetMySQLInstance(nil)

	// v1 路由分组
	v1 := g.Group("/v1")
	{
		userController := user.New(storeIns)

		// users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.POST("", userController.Create) // 创建用户
		}
	}
}
