// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package util

import (
	"context"
	"fmt"
	"sync"

	"github.com/changaolee/skeleton/internal/pkg/scheme"
	"github.com/changaolee/skeleton/pkg/runtime"
	"github.com/spf13/pflag"

	"github.com/changaolee/skeleton/internal/pkg/clioptions"

	"github.com/changaolee/skeleton/internal/pkg/client"
	"github.com/changaolee/skeleton/internal/pkg/rest"
	"github.com/changaolee/skeleton/pkg/version"
)

const (
	flagMatchBinaryVersion = "match-server-version"
)

type MatchVersionFlags struct {
	Delegate clioptions.RESTClientGetter

	RequireMatchedServerVersion bool
	checkServerVersion          sync.Once
	matchesServerVersionErr     error
}

var _ clioptions.RESTClientGetter = &MatchVersionFlags{}

func (f *MatchVersionFlags) checkMatchingServerVersion() error {
	f.checkServerVersion.Do(func() {
		if !f.RequireMatchedServerVersion {
			return
		}

		clientConfig, err := f.Delegate.ToRESTConfig()
		if err != nil {
			f.matchesServerVersionErr = err
			return
		}

		_ = setSKTDefaults(clientConfig)
		restClient, err := rest.RESTClientFor(clientConfig)
		if err != nil {
			f.matchesServerVersionErr = err
			return
		}

		var sVer *version.Info
		if err := restClient.Get().AbsPath("/version").Do(context.TODO()).Into(&sVer); err != nil {
			f.matchesServerVersionErr = err
			return
		}

		clientVersion := version.Get()

		// GitVersion includes GitCommit and GitTreeState, but best to be safe?
		if clientVersion.GitVersion != sVer.GitVersion || clientVersion.GitCommit != sVer.GitCommit ||
			clientVersion.GitTreeState != sVer.GitTreeState {
			f.matchesServerVersionErr = fmt.Errorf(
				"server version (%#v) differs from client version (%#v)",
				sVer,
				version.Get(),
			)
		}
	})

	return f.matchesServerVersionErr
}

func (f *MatchVersionFlags) ToRESTConfig() (*rest.Config, error) {
	if err := f.checkMatchingServerVersion(); err != nil {
		return nil, err
	}
	clientConfig, err := f.Delegate.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	_ = setSKTDefaults(clientConfig)
	return clientConfig, nil
}

func (f *MatchVersionFlags) ToRawSKTConfigLoader() client.ClientConfig {
	return f.Delegate.ToRawSKTConfigLoader()
}

func (f *MatchVersionFlags) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(
		&f.RequireMatchedServerVersion,
		flagMatchBinaryVersion,
		f.RequireMatchedServerVersion,
		"Require server version to match client version",
	)
}

func NewMatchVersionFlags(delegate clioptions.RESTClientGetter) *MatchVersionFlags {
	return &MatchVersionFlags{
		Delegate: delegate,
	}
}

func setSKTDefaults(config *rest.Config) error {
	config.GroupVersion = &scheme.GroupVersion{Group: "skt.api", Version: "v1"}

	if config.APIPath == "" {
		config.APIPath = "/api"
	}
	if config.Negotiator == nil {
		config.Negotiator = runtime.NewSimpleClientNegotiator()
	}
	return rest.SetSKTDefaults(config)
}
