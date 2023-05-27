// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package util

import (
	"github.com/changaolee/skeleton/internal/pkg/client"
	"github.com/changaolee/skeleton/internal/pkg/clioptions"
	"github.com/changaolee/skeleton/internal/pkg/rest"
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
	_ = setSKTDefaults(clientConfig)
	return rest.RESTClientFor(clientConfig)
}
