// Copyright 2022 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package skeleton

import (
	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/internal/pkg/errno"
	"github.com/changaolee/skeleton/internal/pkg/log"
	mw "github.com/changaolee/skeleton/internal/pkg/middleware"
	"github.com/changaolee/skeleton/internal/skeleton/controller/v1/user"
	"github.com/changaolee/skeleton/internal/skeleton/store"
	"github.com/gin-gonic/gin"
)

// installRouters 安装 skeleton 接口路由
func installRouters(g *gin.Engine) error {
	// 注册 404 Handler
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	uc := user.New(store.S)

	g.POST("/login", uc.Login)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.Use(mw.Authn())
			{
				userv1.POST("", uc.Create)
				userv1.PUT(":name/change-password", uc.ChangePassword)
			}
		}
	}

	return nil
}
