// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package rest

import (
	"net/url"
	"strings"

	"github.com/changaolee/skeleton/internal/pkg/scheme"
	"github.com/changaolee/skeleton/pkg/runtime"
	"github.com/changaolee/skeleton/third_party/gorequest"
)

// Interface 定义了和 SKT API 的全部交互方法.
type Interface interface {
	Verb(verb string) *Request
	Post() *Request
	Put() *Request
	Get() *Request
	Delete() *Request
	APIVersion() scheme.GroupVersion
}

// RESTClient 是执行通用 REST 请求的 Client.
type RESTClient struct {
	base             *url.URL              // base 是所有客户端调用的根 URL
	group            string                // group 代表客户端分组，如：skt.api, skt.authz
	versionedAPIPath string                // versionedAPIPath 是定位资源的路径段
	content          ClientContentConfig   // content 描述了如何对响应进行编解码
	Client           *gorequest.SuperAgent // Client 是通用的请求对象
}

// NewRESTClient 创建一个新的 RESTClient.
func NewRESTClient(baseURL *url.URL, versionedAPIPath string,
	config ClientContentConfig, client *gorequest.SuperAgent,
) (*RESTClient, error) {
	if len(config.ContentType) == 0 {
		config.ContentType = "application/json"
	}

	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}

	base.RawQuery = ""
	base.Fragment = ""

	return &RESTClient{
		base:             &base,
		group:            config.GroupVersion.Group,
		versionedAPIPath: versionedAPIPath,
		content:          config,
		Client:           client,
	}, nil
}

func (c *RESTClient) Verb(verb string) *Request {
	return NewRequest(c).Verb(verb)
}

func (c *RESTClient) Post() *Request {
	return c.Verb("POST")
}

func (c *RESTClient) Put() *Request {
	return c.Verb("PUT")
}

func (c *RESTClient) Get() *Request {
	return c.Verb("GET")
}

func (c *RESTClient) Delete() *Request {
	return c.Verb("DELETE")
}

// APIVersion 返回当前 RESTClient 期望使用的 API 版本.
func (c *RESTClient) APIVersion() scheme.GroupVersion {
	return c.content.GroupVersion
}

// ClientContentConfig 控制 RESTClient 与服务器的通信方式.
type ClientContentConfig struct {
	Username string
	Password string

	SecretID  string
	SecretKey string

	BearerToken     string
	BearerTokenFile string

	TLSClientConfig

	AcceptContentTypes string
	ContentType        string
	GroupVersion       scheme.GroupVersion
	Negotiator         runtime.ClientNegotiator
}

// HasBasicAuth 返回配置是否具有 Basic 身份验证.
func (c *ClientContentConfig) HasBasicAuth() bool {
	return len(c.Username) != 0
}

// HasTokenAuth 返回配置是否具有 Token 身份验证.
func (c *ClientContentConfig) HasTokenAuth() bool {
	return len(c.BearerToken) != 0 || len(c.BearerTokenFile) != 0
}

// HasKeyAuth 返回配置是否具有 secretId/secretKey 身份验证.
func (c *ClientContentConfig) HasKeyAuth() bool {
	return len(c.SecretID) != 0 && len(c.SecretKey) != 0
}
