package apiserver

import (
	"github.com/changaolee/skeleton/internal/apiserver/config"
	genericapiserver "github.com/changaolee/skeleton/internal/pkg/server"
	"github.com/changaolee/skeleton/pkg/shutdown"
	"github.com/changaolee/skeleton/pkg/shutdown/managers/posixsignal"
)

type apiServer struct {
	gs               *shutdown.GracefulShutdown
	genericAPIServer *genericapiserver.GenericAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	// 新建优雅关闭组件
	gs := shutdown.New()
	gs.AddManager(posixsignal.NewPosixSignalManager())

	// APIServer
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	// todo: gRPCServer

	server := &apiServer{
		gs:               gs,
		genericAPIServer: genericServer,
	}

	return server, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, err error) {
	genericConfig = genericapiserver.NewConfig()

	// 将 cfg 中的配置更新到 genericConfig 中
	if err = cfg.GenericServerRunOptions.ApplyTo(genericConfig); err != nil {
		return
	}
	if err = cfg.InsecureServing.ApplyTo(genericConfig); err != nil {
		return
	}
	if err = cfg.SecureServing.ApplyTo(genericConfig); err != nil {
		return
	}

	return
}

func (s *apiServer) PrepareRun() *preparedAPIServer {
	initRouter(s.genericAPIServer.Engine)

	s.gs.AddCallback(shutdown.CallbackFunc(func(string) error {
		// todo: mysql、grpc 优雅关闭

		s.genericAPIServer.Shutdown()

		return nil
	}))

	return &preparedAPIServer{s}
}

func (s *preparedAPIServer) Run() error {
	return s.genericAPIServer.Run()
}
