package auth

import (
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/changaolee/skeleton/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// AuthzAudience 定义 jwt 中 audience 字段的值.
const AuthzAudience = "skt.authz.changaolee.com"

type JWTStrategy struct {
	ginjwt.GinJWTMiddleware
}

var _ middleware.AuthStrategy = &JWTStrategy{}

// NewJWTStrategy 创建一个 JWT Bearer 认证策略.
func NewJWTStrategy(gjwt ginjwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{gjwt}
}

// AuthFunc 定义 JWT Bearer 认证策略作为 Gin 中间件.
func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return j.MiddlewareFunc()
}
