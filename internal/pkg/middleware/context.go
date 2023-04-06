package middleware

import (
	"github.com/changaolee/skeleton/pkg/log"
	"github.com/gin-gonic/gin"
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
