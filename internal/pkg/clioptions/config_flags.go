// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package clioptions

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/changaolee/skeleton/internal/pkg/client"
	"github.com/changaolee/skeleton/internal/pkg/rest"
)

// 定义 sktctl 的标志.
const (
	FlagSKTConfig     = "sktconfig"
	FlagBearerToken   = "user.token"
	FlagUsername      = "user.username"
	FlagPassword      = "user.password"
	FlagSecretID      = "user.secret-id"
	FlagSecretKey     = "user.secret-key"
	FlagCertFile      = "user.client-certificate"
	FlagKeyFile       = "user.client-key"
	FlagTLSServerName = "server.tls-server-name"
	FlagInsecure      = "server.insecure-skip-tls-verify"
	FlagCAFile        = "server.certificate-authority"
	FlagAPIServer     = "server.address"
	FlagTimeout       = "server.timeout"
	FlagMaxRetries    = "server.max-retries"
	FlagRetryInterval = "server.retry-interval"
)

type ConfigFlags struct {
	SKTConfig *string

	BearerToken *string
	Username    *string
	Password    *string
	SecretID    *string
	SecretKey   *string

	Insecure      *bool
	TLSServerName *string
	CertFile      *string
	KeyFile       *string
	CAFile        *string

	APIServer     *string
	Timeout       *time.Duration
	MaxRetries    *int
	RetryInterval *time.Duration

	clientConfig client.ClientConfig
	lock         sync.Mutex

	usePersistentConfig bool
}

type RESTClientGetter interface {
	ToRESTConfig() (*rest.Config, error)
	ToRawSKTConfigLoader() client.ClientConfig
}

var _ RESTClientGetter = &ConfigFlags{}

func (f *ConfigFlags) ToRESTConfig() (*rest.Config, error) {
	return f.ToRawSKTConfigLoader().ClientConfig()
}

func (f *ConfigFlags) ToRawSKTConfigLoader() client.ClientConfig {
	if f.usePersistentConfig {
		return f.toRawSKTPersistentConfigLoader()
	}
	return f.toRawSKTConfigLoader()
}

func (f *ConfigFlags) toRawSKTPersistentConfigLoader() client.ClientConfig {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.clientConfig == nil {
		f.clientConfig = f.toRawSKTConfigLoader()
	}

	return f.clientConfig
}

func (f *ConfigFlags) toRawSKTConfigLoader() client.ClientConfig {
	config := client.NewConfig()
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return client.NewClientConfigFromConfig(config)
}

// AddFlags 将客户端配置标志绑定到给定的 FlagSet.
func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.SKTConfig != nil {
		flags.StringVar(
			f.SKTConfig,
			FlagSKTConfig,
			*f.SKTConfig,
			fmt.Sprintf("Path to the %s file to use for CLI requests", FlagSKTConfig),
		)
	}

	if f.BearerToken != nil {
		flags.StringVar(
			f.BearerToken,
			FlagBearerToken,
			*f.BearerToken,
			"Bearer token for authentication to the API server",
		)
	}

	if f.Username != nil {
		flags.StringVar(
			f.Username,
			FlagUsername,
			*f.Username,
			"Username for basic authentication to the API server",
		)
	}

	if f.Password != nil {
		flags.StringVar(
			f.Password,
			FlagPassword,
			*f.Password,
			"Password for basic authentication to the API server",
		)
	}

	if f.SecretID != nil {
		flags.StringVar(
			f.SecretID,
			FlagSecretID,
			*f.SecretID,
			"SecretID for JWT authentication to the API server",
		)
	}

	if f.SecretKey != nil {
		flags.StringVar(
			f.SecretKey,
			FlagSecretKey,
			*f.SecretKey,
			"SecretKey for jwt authentication to the API server",
		)
	}

	if f.CertFile != nil {
		flags.StringVar(
			f.CertFile,
			FlagCertFile,
			*f.CertFile,
			"Path to a client certificate file for TLS",
		)
	}
	if f.KeyFile != nil {
		flags.StringVar(
			f.KeyFile,
			FlagKeyFile,
			*f.KeyFile,
			"Path to a client key file for TLS",
		)
	}
	if f.TLSServerName != nil {
		flags.StringVar(
			f.TLSServerName,
			FlagTLSServerName,
			*f.TLSServerName,
			"Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used",
		)
	}
	if f.Insecure != nil {
		flags.BoolVar(
			f.Insecure,
			FlagInsecure,
			*f.Insecure,
			"If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure",
		)
	}
	if f.CAFile != nil {
		flags.StringVar(
			f.CAFile,
			FlagCAFile,
			*f.CAFile,
			"Path to a cert file for the certificate authority",
		)
	}

	if f.APIServer != nil {
		flags.StringVarP(
			f.APIServer,
			FlagAPIServer,
			"s",
			*f.APIServer,
			"The address and port of the SKT API server",
		)
	}

	if f.Timeout != nil {
		flags.DurationVar(
			f.Timeout,
			FlagTimeout,
			*f.Timeout,
			"The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.",
		)
	}

	if f.MaxRetries != nil {
		flag.IntVar(
			f.MaxRetries,
			FlagMaxRetries,
			*f.MaxRetries,
			"Maximum number of retries.",
		)
	}

	if f.RetryInterval != nil {
		flags.DurationVar(
			f.RetryInterval,
			FlagRetryInterval,
			*f.RetryInterval,
			"The interval time between each attempt.",
		)
	}
}

// WithDeprecatedPasswordFlag 启用 username 和 password 标志.
func (f *ConfigFlags) WithDeprecatedPasswordFlag() *ConfigFlags {
	f.Username = pointer.ToString("")
	f.Password = pointer.ToString("")

	return f
}

// WithDeprecatedSecretFlag 启用 secretID 和 secretKey 标志.
func (f *ConfigFlags) WithDeprecatedSecretFlag() *ConfigFlags {
	f.SecretID = pointer.ToString("")
	f.SecretKey = pointer.ToString("")

	return f
}

// NewConfigFlags 返回设置了默认值的 ConfigFlags.
func NewConfigFlags(usePersistentConfig bool) *ConfigFlags {
	return &ConfigFlags{
		SKTConfig: pointer.ToString(""),

		BearerToken:   pointer.ToString(""),
		Insecure:      pointer.ToBool(false),
		TLSServerName: pointer.ToString(""),
		CertFile:      pointer.ToString(""),
		KeyFile:       pointer.ToString(""),
		CAFile:        pointer.ToString(""),

		APIServer:           pointer.ToString(""),
		Timeout:             pointer.ToDuration(30 * time.Second),
		MaxRetries:          pointer.ToInt(0),
		RetryInterval:       pointer.ToDuration(1 * time.Second),
		usePersistentConfig: usePersistentConfig,
	}
}
