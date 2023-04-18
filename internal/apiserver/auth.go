// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package apiserver

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/changaolee/skeleton/internal/apiserver/store"
	"github.com/changaolee/skeleton/internal/pkg/middleware"
	"github.com/changaolee/skeleton/internal/pkg/middleware/auth"
	"github.com/changaolee/skeleton/internal/pkg/model"
	"github.com/changaolee/skeleton/pkg/log"
)

const (
	// APIServerAudience 定义了 JWT audience 字段的值.
	APIServerAudience = "skt.api.changaolee.com"

	// APIServerIssuer 定义了 JWT issuer 字段的值.
	APIServerIssuer = "skt-apiserver"
)

type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required,username"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

func newAutoAuth() middleware.AuthStrategy {
	return auth.NewAutoStrategy(
		newBasicAuth().(auth.BasicStrategy),
		newJWTAuth().(auth.JWTStrategy),
	)
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username string, password string) bool {
		// 从数据库中拉取用户信息
		user, err := store.Store().Users().Get(context.TODO(), username)
		if err != nil {
			return false
		}

		// 验证用户密码
		if err := user.Compare(password); err != nil {
			return false
		}

		user.LoginAt = time.Now()
		_ = store.Store().Users().Update(context.TODO(), user)

		return true
	})
}

func newJWTAuth() middleware.AuthStrategy {
	ginjwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.realm"), // JWT 标识
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.key")),   // 用于签名的密钥
		Timeout:          viper.GetDuration("jwt.timeout"),     // token 有效时间
		MaxRefresh:       viper.GetDuration("jwt.max-refresh"), // token 最长更新间隔
		Authenticator:    authenticator(),                      // 用户身份验证，
		LoginResponse:    loginResponse(),                      // 登录响应
		LogoutResponse: func(c *gin.Context, code int) { // 退出登录响应
			c.JSON(http.StatusOK, nil)
		},
		RefreshResponse: refreshResponse(), // token 刷新响应
		PayloadFunc:     payloadFunc(),     // 添加额外业务相关的信息
		IdentityHandler: func(c *gin.Context) interface{} { // 提取身份标识
			claims := jwt.ExtractClaims(c)
			return claims[jwt.IdentityKey]
		},
		IdentityKey:  middleware.UsernameKey, // 身份标识字段
		Authorizator: authorizator(),         // 用户权限检查
		Unauthorized: func(c *gin.Context, code int, message string) { // 未授权响应
			c.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,     // 是否将 token 作为 cookie 返回
		TimeFunc:      time.Now, // 提供当前时间
	})
	return auth.NewJWTStrategy(*ginjwt)
}

// authenticator 验证用户身份.
func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login loginInfo
		var err error

		// 支持从 header 或 body 中提取
		if c.Request.Header.Get("Authorization") != "" {
			login, err = parseWithHeader(c)
		} else {
			login, err = parseWithBody(c)
		}
		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		// 基于登录用户名查找用户
		user, err := store.Store().Users().Get(c, login.Username)
		if err != nil {
			log.Errorf("Get user information failed: %s", err.Error())

			return "", jwt.ErrFailedAuthentication
		}

		// 验证用户密码
		if err := user.Compare(login.Password); err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		user.LoginAt = time.Now()
		_ = store.Store().Users().Update(c, user)

		return user, nil
	}
}

func parseWithHeader(c *gin.Context) (loginInfo, error) {
	authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(authHeader) != 2 || authHeader[0] != "Basic" {
		log.Errorf("Get basic string from Authorization header failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	payload, err := base64.StdEncoding.DecodeString(authHeader[1])
	if err != nil {
		log.Errorf("Decode basic string: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		log.Errorf("Parse payload failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return loginInfo{
		Username: pair[0],
		Password: pair[1],
	}, nil
}

func parseWithBody(c *gin.Context) (loginInfo, error) {
	var login loginInfo
	if err := c.ShouldBindJSON(&login); err != nil {
		log.Errorf("Parse login parameters: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}
	return login, nil
}

func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		claims := jwt.MapClaims{
			"iss": APIServerIssuer,
			"aud": APIServerAudience,
		}
		if u, ok := data.(*model.User); ok {
			claims[jwt.IdentityKey] = u.Name
			claims["sub"] = u.Name
		}
		return claims
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(string); ok {
			log.C(c).Infof("User `%s` is authenticated.", v)
			return true
		}
		return false
	}
}
