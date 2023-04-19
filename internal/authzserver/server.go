package authzserver

import (
	"context"

	"github.com/changaolee/skeleton/internal/authzserver/cache"
	"github.com/changaolee/skeleton/internal/authzserver/config"
	"github.com/changaolee/skeleton/internal/authzserver/load"
	"github.com/changaolee/skeleton/internal/authzserver/store/apiserver"
	genericapiserver "github.com/changaolee/skeleton/internal/pkg/server"
	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/changaolee/skeleton/pkg/shutdown"
	"github.com/changaolee/skeleton/pkg/shutdown/managers/posixsignal"
)

type authzServer struct {
	rpcServer        string
	clientCA         string
	gs               *shutdown.GracefulShutdown
	genericAPIServer *genericapiserver.GenericAPIServer
	redisCancelFunc  context.CancelFunc
}

type preparedAuthzServer struct {
	*authzServer
}

func createAuthzServer(cfg *config.Config) (*authzServer, error) {
	// 优雅关闭组件
	gs := shutdown.New()
	gs.AddManager(posixsignal.NewPosixSignalManager())

	// Redis 实例
	_, err := cache.GetRedisInstance(cfg.RedisOptions)
	if err != nil {
		return nil, err
	}

	// APIServer
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &authzServer{
		rpcServer:        cfg.RPCServer,
		clientCA:         cfg.ClientCA,
		gs:               gs,
		genericAPIServer: genericServer,
	}

	return server, nil
}

func (s *authzServer) PrepareRun() *preparedAuthzServer {
	_ = s.initialize()

	initRouter(s.genericAPIServer.Engine)

	s.gs.AddCallback(shutdown.CallbackFunc(func(string) error {
		s.genericAPIServer.Shutdown()
		s.redisCancelFunc()

		return nil
	}))

	return &preparedAuthzServer{s}
}

func (s *authzServer) initialize() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.redisCancelFunc = cancel

	// 后台定时从 apiserver 重载 secrets 和 policies.
	client, err := apiserver.GetAPIServerCacheClientInstance(s.rpcServer, s.clientCA)
	if err != nil {
		return errors.Wrap(err, "get cache client failed")
	}
	cacheIns, err := load.GetCacheInstance(client)
	if err != nil {
		return errors.Wrap(err, "get cache instance failed")
	}
	load.NewLoader(ctx, cacheIns).Start()

	// todo: analytics

	return nil
}

func (s *preparedAuthzServer) Run() error {
	return s.genericAPIServer.Run()
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, err error) {
	genericConfig = genericapiserver.NewConfig()

	// 将 cfg 中的配置更新到 genericConfig 中
	if err = cfg.GenericServerRunOptions.ApplyTo(genericConfig); err != nil {
		return
	}
	if err = cfg.SecureServing.ApplyTo(genericConfig); err != nil {
		return
	}
	if err = cfg.InsecureServing.ApplyTo(genericConfig); err != nil {
		return
	}

	return
}
