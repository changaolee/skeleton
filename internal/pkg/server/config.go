package server

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	// todo: 配置路由、中间件

	return s, nil
}
