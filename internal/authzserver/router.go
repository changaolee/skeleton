package authzserver

import (
	"github.com/changaolee/skeleton/internal/authzserver/cache"
	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/changaolee/skeleton/pkg/log"
	"github.com/gin-gonic/gin"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	auth := newCacheAuth()
	g.NoRoute(auth.AuthFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "page not found."), nil)
	})

	cacheIns, _ := cache.GetRedisInstance(nil)
	if cacheIns == nil {
		log.Panicf("Get nil cache instance")
	}

	// todo: authz 服务路由注册
	//apiv1 := g.Group("/v1", auth.AuthFunc())
	//{
	//	authzController := authorize.NewAuthzController(cacheIns)
	//
	//	// Router for authorization
	//	apiv1.POST("/authz", authzController.Authorize)
	//}

	return g
}
