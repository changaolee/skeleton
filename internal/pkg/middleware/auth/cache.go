package auth

import (
	"fmt"
	"time"

	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/internal/pkg/middleware"
	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrMissingKID    = errors.New("Invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("Can not obtain secret information from cache")
)

type Secret struct {
	Username string
	ID       string
	Key      string
	Expires  int64
}

type getSecretFunc func(kid string) (Secret, error)

// CacheStrategy 定义 Cache 认证策略（基于缓存实现的 JWT Bearer 认证）.
type CacheStrategy struct {
	get getSecretFunc
}

var _ middleware.AuthStrategy = &CacheStrategy{}

// NewCacheStrategy 基于给定的 get 方法创建一个 Cache 认证策略.
func NewCacheStrategy(get getSecretFunc) CacheStrategy {
	return CacheStrategy{get}
}

func (cache CacheStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if len(header) == 0 {
			core.WriteResponse(c, errors.WithCode(code.ErrMissingHeader, "Authorization header cannot be empty."), nil)
			c.Abort()

			return
		}

		var rawJWT string
		_, _ = fmt.Sscanf(header, "Bearer %s", &rawJWT)

		// 下面使用自定义的验证逻辑
		var secret Secret
		claims := &jwt.MapClaims{}

		// 验证 token
		parsedT, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (interface{}, error) {
			// 验证 token 的签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// 从 token 中获取 secret 的标识 kid
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, ErrMissingKID
			}

			// 从缓存中读取 secret
			secret, err := cache.get(kid)
			if err != nil {
				return nil, ErrMissingSecret
			}

			// 返回 secret
			return []byte(secret.Key), nil
		})
		// 在 ParseWithClaims 时会自动为 parsedT 的 Valid 赋值
		if err != nil || !parsedT.Valid {
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, err.Error()), nil)
			c.Abort()

			return
		}

		if KeyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			core.WriteResponse(c, errors.WithCode(code.ErrExpired, "expired at: %s", tm), nil)
			c.Abort()

			return
		}

		c.Set(middleware.UsernameKey, secret.Username)
		c.Next()
	}
}

// KeyExpired 检查一个 key 是否过期.
func KeyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0))
	}
	return false
}
