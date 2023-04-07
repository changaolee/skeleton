package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/internal/pkg/middleware"
	"github.com/changaolee/skeleton/pkg/log"
	"github.com/changaolee/skeleton/pkg/version"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// GenericAPIServer 包含一个 skeleton api 服务器的配置.
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

	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

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

// ping 对服务进行健康检查.
func (s *GenericAPIServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)
	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServingInfo.Address, ":")[1])
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Infow("The router has been deployed successfully.")

			_ = resp.Body.Close()

			return nil
		}

		log.Infow("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatalw("Can not ping http server within the specified time interval.")
			return err
		default:
		}
	}
}

func initGenericAPIServer(s *GenericAPIServer) {
	s.setup()
	s.installMiddlewares()
	s.installAPIs()
}

// setup 针对 gin 引擎做一些设置.
func (s *GenericAPIServer) setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// installMiddlewares 添加中间件.
func (s *GenericAPIServer) installMiddlewares() {
	// 必要的中间件
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())

	// 自定义中间件
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("Can not find middleware: %s", m)
			continue
		}
		log.Infof("Install middleware: %s", m)
		s.Use(mw)
	}
}

// installAPIs 添加 api.
func (s *GenericAPIServer) installAPIs() {
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			core.WriteResponse(c, nil, map[string]string{"status": "ok"})
		})
	}

	s.GET("/version", func(c *gin.Context) {
		core.WriteResponse(c, nil, version.Get())
	})
}
