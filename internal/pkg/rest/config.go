// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package rest

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/changaolee/skeleton/internal/pkg/scheme"
	"github.com/changaolee/skeleton/pkg/runtime"
	"github.com/changaolee/skeleton/pkg/third_party/gorequest"
)

// Config 包含可以在初始化时传递给 SKT 客户端的通用属性.
type Config struct {
	Host    string
	APIPath string
	ContentConfig

	// Server requires Basic authentication
	Username string
	Password string

	SecretID  string
	SecretKey string

	// Server requires Bearer authentication. This client will not attempt to use
	// refresh tokens for an OAuth2 flow.
	// TODO: demonstrate an OAuth2 compatible client.
	BearerToken string

	// Path to a file containing a BearerToken.
	// If set, the contents are periodically read.
	// The last successfully read value takes precedence over BearerToken.
	BearerTokenFile string

	// TLSClientConfig contains settings to enable transport layer security
	TLSClientConfig

	// UserAgent is an optional field that specifies the caller of this request.
	UserAgent string
	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout       time.Duration
	MaxRetries    int
	RetryInterval time.Duration
}

type ContentConfig struct {
	ServiceName        string
	AcceptContentTypes string
	ContentType        string
	GroupVersion       *scheme.GroupVersion
	Negotiator         runtime.ClientNegotiator
}

// TLSClientConfig 包含启用传输层安全的设置.
type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	Insecure bool
	// ServerName is passed to the server for SNI and is used in the client to check server
	// ceritificates against. If ServerName is empty, the hostname used to contact the
	// server is used.
	ServerName string

	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string

	// CertData holds PEM-encoded bytes (typically read from a client certificate file).
	// CertData takes precedence over CertFile
	CertData []byte
	// KeyData holds PEM-encoded bytes (typically read from a client certificate key file).
	// KeyData takes precedence over KeyFile
	KeyData []byte
	// CAData holds PEM-encoded bytes (typically read from a root certificates bundle).
	// CAData takes precedence over CAFile
	CAData []byte

	// NextProtos is a list of supported application level protocols, in order of preference.
	// Used to populate tls.Config.NextProtos.
	// To indicate to the server http/1.1 is preferred over http/2, set to ["http/1.1", "h2"] (though the server is free
	// to ignore that preference).
	// To use only http/1.1, set to ["http/1.1"].
	NextProtos []string
}

func RESTClientFor(config *Config) (*RESTClient, error) {
	if config.GroupVersion == nil {
		return nil, fmt.Errorf("GroupVersion is required when initializing a RESTClient")
	}

	if config.Negotiator == nil {
		return nil, fmt.Errorf("NegotiatedSerializer is required when initializing a RESTClient")
	}

	baseURL, versionedAPIPath, err := defaultServerURLFor(config)
	if err != nil {
		return nil, err
	}

	// Get the TLS options for this client config
	tlsConfig, err := TLSConfigFor(config)
	if err != nil {
		return nil, err
	}

	// Only retry when get a server side error.
	client := gorequest.New().TLSClientConfig(tlsConfig).Timeout(config.Timeout).
		Retry(config.MaxRetries, config.RetryInterval, http.StatusInternalServerError)
	// NOTICE: must set DoNotClearSuperAgent to true, or the client will clean header befor http.Do
	client.DoNotClearSuperAgent = true

	var gv scheme.GroupVersion
	if config.GroupVersion != nil {
		gv = *config.GroupVersion
	}

	clientContent := ClientContentConfig{
		Username:           config.Username,
		Password:           config.Password,
		SecretID:           config.SecretID,
		SecretKey:          config.SecretKey,
		BearerToken:        config.BearerToken,
		BearerTokenFile:    config.BearerTokenFile,
		TLSClientConfig:    config.TLSClientConfig,
		AcceptContentTypes: config.AcceptContentTypes,
		ContentType:        config.ContentType,
		GroupVersion:       gv,
		Negotiator:         config.Negotiator,
	}

	return NewRESTClient(baseURL, versionedAPIPath, clientContent, client)
}

func TLSConfigFor(c *Config) (*tls.Config, error) {
	if !(c.HasCA() || c.HasCertAuth() || c.Insecure || len(c.ServerName) > 0) {
		return nil, nil
	}

	if c.HasCA() && c.Insecure {
		return nil, fmt.Errorf("specifying a root certificates file with the insecure flag is not allowed")
	}

	if err := LoadTLSFiles(c); err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		// Can't use SSLv3 because of POODLE and BEAST
		// Can't use TLSv1.0 because of POODLE and BEAST using CBC cipher
		// Can't use TLSv1.1 because of RC4 cipher usage
		MinVersion: tls.VersionTLS12,

		InsecureSkipVerify: c.Insecure,
		ServerName:         c.ServerName,
		NextProtos:         c.NextProtos,
	}

	if c.HasCA() {
		tlsConfig.RootCAs = rootCertPool(c.CAData)
	}

	var staticCert *tls.Certificate
	// Treat cert as static if either key or cert was data, not a file
	if c.HasCertAuth() {
		// If key/cert were provided, verify them before setting up
		// tlsConfig.GetClientCertificate.
		cert, err := tls.X509KeyPair(c.CertData, c.KeyData)
		if err != nil {
			return nil, err
		}

		staticCert = &cert
	}

	if c.HasCertAuth() {
		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			// Note: static key/cert data always take precedence over cert
			// callback.
			if staticCert != nil {
				return staticCert, nil
			}

			// Both c.TLS.CertData/KeyData were unset and GetCert didn't return
			// anything. Return an empty tls.Certificate, no client cert will
			// be sent to the server.
			return &tls.Certificate{}, nil
		}
	}

	return tlsConfig, nil
}

// HasCA 返回配置是否具有证书颁发机构.
func (c TLSClientConfig) HasCA() bool {
	return len(c.CAData) > 0 || len(c.CAFile) > 0
}

// HasCertAuth 返回配置是否具有证书身份验证.
func (c TLSClientConfig) HasCertAuth() bool {
	return (len(c.CertData) != 0 || len(c.CertFile) != 0) && (len(c.KeyData) != 0 || len(c.KeyFile) != 0)
}

// LoadTLSFiles 将数据从 CertFile、KeyFile 和 CAFile 复制到 CertData、KeyData 和 CAFile 中，或返回错误.
func LoadTLSFiles(c *Config) error {
	var err error

	c.CAData, err = dataFromSliceOrFile(c.CAData, c.CAFile)
	if err != nil {
		return err
	}

	c.CertData, err = dataFromSliceOrFile(c.CertData, c.CertFile)
	if err != nil {
		return err
	}

	c.KeyData, err = dataFromSliceOrFile(c.KeyData, c.KeyFile)
	if err != nil {
		return err
	}

	return nil
}

func dataFromSliceOrFile(data []byte, file string) ([]byte, error) {
	if len(data) > 0 {
		return base64.StdEncoding.DecodeString(string(data))
	}

	if len(file) > 0 {
		fileData, err := os.ReadFile(file)
		if err != nil {
			return []byte{}, err
		}

		return fileData, nil
	}

	return nil, nil
}

func rootCertPool(caData []byte) *x509.CertPool {
	if len(caData) == 0 {
		return nil
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)

	return certPool
}
