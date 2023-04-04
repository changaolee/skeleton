package apiserver

import (
	"github.com/changaolee/skeleton/internal/apiserver/config"
	genericapiserver "github.com/changaolee/skeleton/internal/pkg/server"
)

type apiServer struct {
	genericAPIServer *genericapiserver.GenericAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
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
		genericAPIServer: genericServer,
	}

	return server, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()

	// todo: 将 cfg 中的配置更新到 genericConfig 中

	return
}

func (s *apiServer) PrepareRun() *preparedAPIServer {
	initRouter(s.genericAPIServer.Engine)
	
	// todo: 优雅关闭

	return &preparedAPIServer{s}
}

func (s *preparedAPIServer) Run() error {
	return s.genericAPIServer.Run()
}
