// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package server

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/changaolee/skeleton/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	// RecommendedHomeDir 定义了所有 skeleton 服务配置文件的默认存放目录.
	RecommendedHomeDir = ".skt"

	// RecommendedEnvPrefix 定义了所有 skeleton 服务环境变量的前缀.
	RecommendedEnvPrefix = "SKT"
)

// Config 是一个通用 API 服务器的配置结构.
type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Jwt             *JwtInfo
	Mode            string
	Middlewares     []string
	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
}

// CertKey 包含与证书相关的配置项.
type CertKey struct {
	// CertFile 是一个包含 PEM 编码的证书的文件，可能还包含完整的证书链.
	CertFile string
	// KeyFile 是一个包含 CertFile 指定证书 PEM 编码私钥的文件.
	KeyFile string
}

// SecureServingInfo 保存 HTTPS 服务配置.
type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

// Address 将 IP 和 port 拼接为一个地址字符串，如：0.0.0.0:8443.
func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// InsecureServingInfo 保存 HTTP 服务配置.
type InsecureServingInfo struct {
	Address string
}

// JwtInfo 定义用于创建 JWT 身份验证中间件的 JWT 字段.
type JwtInfo struct {
	Realm      string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
}

// NewConfig 创建一个默认 Config.
func NewConfig() *Config {
	return &Config{
		Jwt: &JwtInfo{
			Realm:      "skt jwt",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		Healthz:         true,
		EnableProfiling: true,
		EnableMetrics:   true,
	}
}

// CompletedConfig 是一个完整的 GenericAPIServer 配置.
type CompletedConfig struct {
	*Config
}

// Complete 基于 Config 补充完整的配置.
func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{c}
}

func (c *CompletedConfig) New() (*GenericAPIServer, error) {
	gin.SetMode(c.Mode)

	s := &GenericAPIServer{
		middlewares:         c.Middlewares,
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		Engine:              gin.New(),
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
	}

	initGenericAPIServer(s)

	return s, nil
}

// LoadConfig 读取配置文件和 ENV 变量.
func LoadConfig(cfg string, defaultName string) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		home, _ := os.UserHomeDir()
		viper.AddConfigPath(filepath.Join(home, RecommendedHomeDir))
		viper.AddConfigPath("/etc/skt")
		viper.SetConfigName(defaultName)
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(RecommendedEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("WARNING: viper failed to discover and load the configuration file: %s", err.Error())
	}
}
