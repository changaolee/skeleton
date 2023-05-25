// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package rest

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/changaolee/skeleton/pkg/auth"
	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/changaolee/skeleton/pkg/third_party/gorequest"
)

type Request struct {
	c *RESTClient

	timeout time.Duration

	verb       string
	pathPrefix string
	subpath    string
	params     url.Values
	headers    http.Header

	resource     string
	resourceName string
	subresource  string

	err  error
	body interface{}
}

// NewRequest 创建一个新的 Request.
func NewRequest(c *RESTClient) *Request {
	var pathPrefix string
	if c.base != nil {
		pathPrefix = path.Join("/", c.base.Path, c.versionedAPIPath)
	} else {
		pathPrefix = path.Join("/", c.versionedAPIPath)
	}

	r := &Request{
		c:          c,
		pathPrefix: pathPrefix,
	}

	authMethod := 0

	for _, fn := range []func() bool{c.content.HasBasicAuth, c.content.HasTokenAuth, c.content.HasKeyAuth} {
		if fn() {
			authMethod++
		}
	}

	if authMethod > 1 {
		r.err = fmt.Errorf(
			"username/password or bearer token or secretID/secretKey may be set, but should use only one of them",
		)

		return r
	}

	switch {
	case c.content.HasTokenAuth():
		r.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.content.BearerToken))
	case c.content.HasKeyAuth():
		tokenString := auth.Sign(c.content.SecretID, c.content.SecretKey, "skeleton", c.group+".changaolee.com")
		r.SetHeader("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	case c.content.HasBasicAuth():
		r.SetHeader("Authorization", "Basic "+basicAuth(c.content.Username, c.content.Password))
	}

	switch {
	case len(c.content.AcceptContentTypes) > 0:
		r.SetHeader("Accept", c.content.AcceptContentTypes)
	case len(c.content.ContentType) > 0:
		r.SetHeader("Accept", c.content.ContentType+", */*")
	}

	return r
}

func basicAuth(username, password string) string {
	ba := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(ba))
}

// Verb 设置 Request 的动作.
func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

// SetHeader 为一个 HTTP Request 设置 Header.
func (r *Request) SetHeader(key string, values ...string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}

	r.headers.Del(key)

	for _, value := range values {
		r.headers.Add(key, value)
	}

	return r
}

// AbsPath 用给定的 segments 重写现有路径.
func (r *Request) AbsPath(segments ...string) *Request {
	if r.err != nil {
		return r
	}

	r.pathPrefix = path.Join(r.c.base.Path, path.Join(segments...))

	if len(segments) == 1 && (len(r.c.base.Path) > 1 || len(segments[0]) > 1) && strings.HasSuffix(segments[0], "/") {
		// preserve any trailing slashes for legacy behavior
		r.pathPrefix += "/"
	}

	return r
}

// URL 返回当前请求的 URL.
func (r *Request) URL() *url.URL {
	p := r.pathPrefix
	if len(r.resource) != 0 {
		p = path.Join(p, strings.ToLower(r.resource))
	}
	if len(r.resourceName) != 0 || len(r.subpath) != 0 || len(r.subresource) != 0 {
		p = path.Join(p, r.resourceName, r.subresource, r.subpath)
	}

	finalURL := &url.URL{}
	if r.c.base != nil {
		*finalURL = *r.c.base
	}

	finalURL.Path = p

	query := url.Values{}

	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	// timeout is handled specially here.
	if r.timeout != 0 {
		query.Set("timeout", r.timeout.String())
	}

	finalURL.RawQuery = query.Encode()

	return finalURL
}

func (r *Request) Do(ctx context.Context) Result {
	client := r.c.Client
	client.Header = r.headers

	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)

		defer cancel()
	}

	client.WithContext(ctx)

	resp, body, errs := client.CustomMethod(r.verb, r.URL().String()).Send(r.body).EndBytes()
	if err := combineErr(resp, body, errs); err != nil {
		return Result{
			response: &resp,
			err:      err,
			body:     body,
		}
	}

	decoder, err := r.c.content.Negotiator.Decoder()
	if err != nil {
		return Result{
			response: &resp,
			err:      err,
			body:     body,
			decoder:  decoder,
		}
	}

	return Result{
		response: &resp,
		body:     body,
		decoder:  decoder,
	}
}

func combineErr(resp gorequest.Response, body []byte, errs []error) error {
	var e, sep string

	if len(errs) > 0 {
		for _, err := range errs {
			e = sep + err.Error()
			sep = "\n"
		}

		return errors.New(e)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}
