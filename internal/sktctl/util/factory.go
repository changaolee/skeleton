package util

import (
	"github.com/changaolee/skeleton/internal/pkg/client"
	"github.com/changaolee/skeleton/internal/pkg/rest"
	"github.com/changaolee/skeleton/pkg/clioptions"
)

type Factory interface {
	clioptions.RESTClientGetter
	RESTClient() (*rest.RESTClient, error)
}

type factory struct {
	clientGetter clioptions.RESTClientGetter
}

func NewFactory(clientGetter clioptions.RESTClientGetter) Factory {
	if clientGetter == nil {
		panic("attempt to instantiate factory with nil clientGetter")
	}

	f := &factory{
		clientGetter: clientGetter,
	}

	return f
}

func (f *factory) ToRESTConfig() (*rest.Config, error) {
	return f.clientGetter.ToRESTConfig()
}

func (f *factory) ToRawSKTConfigLoader() client.ClientConfig {
	return f.clientGetter.ToRawSKTConfigLoader()
}

func (f *factory) RESTClient() (*rest.RESTClient, error) {
	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return rest.RESTClientFor(clientConfig)
}
