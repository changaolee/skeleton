package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/changaolee/skeleton/pkg/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// GenericAPIServer contains state for an iam api server.
// type GenericAPIServer gin.Engine.
type GenericAPIServer struct {
	middlewares []string

	// SecureServingInfo 保存 HTTPS 服务配置
	SecureServingInfo *SecureServingInfo

	// InsecureServingInfo 保存 HTTP 服务配置
	InsecureServingInfo *InsecureServingInfo

	// ShutdownTimeout 表示优雅关闭的超时时间
	ShutdownTimeout time.Duration

	*gin.Engine
	healthz         bool
	enableMetrics   bool
	enableProfiling bool

	insecureServer, secureServer *http.Server
}

func (s *GenericAPIServer) Run() error {
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
	}

	s.secureServer = &http.Server{
		Addr:    s.SecureServingInfo.Address(),
		Handler: s,
	}

	var eg errgroup.Group

	eg.Go(func() error {
		log.Infof("Start to listening the incoming requests on http address: %s", s.InsecureServingInfo.Address)

		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	eg.Go(func() error {
		key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
		if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
			return nil
		}

		log.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())

		if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.SecureServingInfo.Address())

		return nil
	})

	// todo: 健康检查

	if err := eg.Wait(); err != nil {
		log.Fatalw(err.Error())
	}

	return nil
}

// Shutdown 优雅关闭 API 服务.
func (s *GenericAPIServer) Shutdown() {
	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.secureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown secure server failed: %s", err.Error())
	}
	if err := s.insecureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown insecure server failed: %s", err.Error())
	}
}
