// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package rest

import (
	"fmt"
	"net/url"
	"path"

	"github.com/changaolee/skeleton/internal/pkg/scheme"
)

func DefaultServerURL(
	host, apiPath string,
	groupVersion scheme.GroupVersion,
	defaultTLS bool,
) (*url.URL, string, error) {
	hostURL, err := url.Parse(host)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		requestURL := fmt.Sprintf("http://%s.changaolee.com:8080", groupVersion.Group)
		if defaultTLS {
			requestURL = fmt.Sprintf("https://%s.changaolee.com:8443", groupVersion.Group)
		}

		hostURL, err = url.Parse(requestURL)
		if err != nil {
			return nil, "", err
		}

		if hostURL.Path != "" && hostURL.Path != "/" {
			return nil, "", fmt.Errorf("host must be a URL or a host:port pair: %q", host)
		}
	}

	versionedAPIPath := path.Join("/", apiPath, groupVersion.Version)

	return hostURL, versionedAPIPath, nil
}

func defaultServerURLFor(config *Config) (*url.URL, string, error) {
	hasCA := len(config.CAFile) != 0 || len(config.CAData) != 0
	hasCert := len(config.CertFile) != 0 || len(config.CertData) != 0
	defaultTLS := hasCA || hasCert || config.Insecure

	if config.GroupVersion != nil {
		return DefaultServerURL(config.Host, config.APIPath, *config.GroupVersion, defaultTLS)
	}

	return DefaultServerURL(config.Host, config.APIPath, scheme.GroupVersion{}, defaultTLS)
}
