// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package options

import (
	"fmt"
	"path"

	"github.com/spf13/pflag"

	"github.com/changaolee/skeleton/internal/pkg/server"
)

// SecureServingOptions 定义了 HTTPS 服务配置.
type SecureServingOptions struct {
	BindAddress string           `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int              `json:"bind-port"    mapstructure:"bind-port"`
	ServerCert  GeneratedKeyCert `json:"tls"          mapstructure:"tls"`
	RequirePort bool
}

// CertKey 包含了证书和私钥.
type CertKey struct {
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`        // 证书文件
	KeyFile  string `json:"private-key-file" mapstructure:"private-key-file"` // 私钥文件
}

// GeneratedKeyCert 包含了证书相关配置.
type GeneratedKeyCert struct {
	CertKey       CertKey `json:"cert-key"  mapstructure:"cert-key"`  // 支持显式指定证书和私钥
	CertDirectory string  `json:"cert-dir"  mapstructure:"cert-dir"`  // 未设置证书和私钥时会在此目录生成
	PairName      string  `json:"pair-name" mapstructure:"pair-name"` // 生成证书和私钥的文件名
}

// NewSecureServingOptions 创建了一个默认 HTTPS 服务配置.
func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8443,
		ServerCert: GeneratedKeyCert{
			PairName:      "skt",
			CertDirectory: "/var/run/skt",
		},
		RequirePort: true,
	}
}

// ApplyTo 将当前选项绑定到 Config 中.
func (s *SecureServingOptions) ApplyTo(c *server.Config) error {
	c.SecureServing = &server.SecureServingInfo{
		BindAddress: s.BindAddress,
		BindPort:    s.BindPort,
		CertKey: server.CertKey{
			CertFile: s.ServerCert.CertKey.CertFile,
			KeyFile:  s.ServerCert.CertKey.KeyFile,
		},
	}

	return nil
}

// Validate 验证 HTTPS 选项.
func (s *SecureServingOptions) Validate() []error {
	if s == nil {
		return nil
	}

	var errors []error

	if s.RequirePort && s.BindPort < 1 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--secure.bind-port %v must be between 1 and 65535, inclusive. It cannot be turned off with 0",
				s.BindPort,
			),
		)
	} else if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(errors, fmt.Errorf("--secure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off secure port", s.BindPort))
	}

	return errors
}

// AddFlags 向指定 FlagSet 中添加 HTTPS 选项相关标志.
func (s *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "secure.bind-address", s.BindAddress, ""+
		"The IP address on which to listen for the --secure.bind-port port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	desc := "The port on which to serve HTTPS with authentication and authorization."
	if s.RequirePort {
		desc += " It cannot be switched off with 0."
	} else {
		desc += " If 0, don't serve HTTPS at all."
	}
	fs.IntVar(&s.BindPort, "secure.bind-port", s.BindPort, desc)

	fs.StringVar(&s.ServerCert.CertDirectory, "secure.tls.cert-dir", s.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&s.ServerCert.PairName, "secure.tls.pair-name", s.ServerCert.PairName, ""+
		"The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&s.ServerCert.CertKey.CertFile, "secure.tls.cert-key.cert-file", s.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "secure.tls.cert-key.private-key-file",
		s.ServerCert.CertKey.KeyFile, ""+
			"File containing the default x509 private key matching --secure.tls.cert-key.cert-file.")
}

// Complete 填写所有需要有效数据但未设置的字段.
func (s *SecureServingOptions) Complete() error {
	if s == nil || s.BindPort == 0 {
		return nil
	}

	keyCert := &s.ServerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 {
		return nil
	}

	if len(s.ServerCert.CertDirectory) > 0 {
		if len(s.ServerCert.PairName) == 0 {
			return fmt.Errorf("--secure.tls.pair-name is required if --secure.tls.cert-dir is set")
		}
		keyCert.CertFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".key")
	}

	return nil
}
