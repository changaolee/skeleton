package client

import (
	"net/url"
	"time"

	"github.com/changaolee/skeleton/internal/pkg/rest"
)

type Server struct {
	LocationOfOrigin         string
	Timeout                  time.Duration `yaml:"timeout,omitempty"                    mapstructure:"timeout,omitempty"`
	MaxRetries               int           `yaml:"max-retries,omitempty"                mapstructure:"max-retries,omitempty"`
	RetryInterval            time.Duration `yaml:"retry-interval,omitempty"             mapstructure:"retry-interval,omitempty"`
	Address                  string        `yaml:"address,omitempty"                    mapstructure:"address,omitempty"`
	TLSServerName            string        `yaml:"tls-server-name,omitempty"            mapstructure:"tls-server-name,omitempty"`            // +optional
	InsecureSkipTLSVerify    bool          `yaml:"insecure-skip-tls-verify,omitempty"   mapstructure:"insecure-skip-tls-verify,omitempty"`   // +optional
	CertificateAuthority     string        `yaml:"certificate-authority,omitempty"      mapstructure:"certificate-authority,omitempty"`      // +optional
	CertificateAuthorityData string        `yaml:"certificate-authority-data,omitempty" mapstructure:"certificate-authority-data,omitempty"` // +optional
}

type AuthInfo struct {
	LocationOfOrigin      string
	ClientCertificate     string `yaml:"client-certificate,omitempty"      mapstructure:"client-certificate,omitempty"`
	ClientCertificateData string `yaml:"client-certificate-data,omitempty" mapstructure:"client-certificate-data,omitempty"` // +optional
	ClientKey             string `yaml:"client-key,omitempty"              mapstructure:"client-key,omitempty"`              // +optional
	ClientKeyData         string `yaml:"client-key-data,omitempty"         mapstructure:"client-key-data,omitempty"`         // +optional
	Token                 string `yaml:"token,omitempty"                   mapstructure:"token,omitempty"`                   // +optional

	Username string `yaml:"username,omitempty" mapstructure:"username,omitempty"`
	Password string `yaml:"password,omitempty" mapstructure:"password,omitempty"`

	SecretID  string `yaml:"secret-id,omitempty"  mapstructure:"secret-id,omitempty"`
	SecretKey string `yaml:"secret-key,omitempty" mapstructure:"secret-key,omitempty"`
}

type Config struct {
	APIVersion string    `yaml:"apiVersion,omitempty" mapstructure:"apiVersion,omitempty"`
	AuthInfo   *AuthInfo `yaml:"user,omitempty"       mapstructure:"user,omitempty"`
	Server     *Server   `yaml:"server,omitempty"     mapstructure:"server,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Server:   &Server{},
		AuthInfo: &AuthInfo{},
	}
}

type ClientConfig interface {
	// ClientConfig returns a complete client config
	ClientConfig() (*rest.Config, error)
}

type DirectClientConfig struct {
	config Config
}

func NewClientConfigFromConfig(config *Config) ClientConfig {
	return &DirectClientConfig{*config}
}

func (config *DirectClientConfig) getAuthInfo() AuthInfo {
	return *config.config.AuthInfo
}

func (config *DirectClientConfig) getServer() Server {
	return *config.config.Server
}

func (config *DirectClientConfig) ConfirmUsable() error {
	validationErrors := make([]error, 0)

	authInfo := config.getAuthInfo()
	validationErrors = append(validationErrors, validateAuthInfo(authInfo)...)
	server := config.getServer()
	validationErrors = append(validationErrors, validateServerInfo(server)...)
	// when direct client config is specified, and our only error is that no server is defined, we should
	// return a standard "no config" error
	if len(validationErrors) == 1 && validationErrors[0] == ErrEmptyServer {
		return newErrConfigurationInvalid([]error{ErrEmptyConfig})
	}

	return newErrConfigurationInvalid(validationErrors)
}

func (config *DirectClientConfig) ClientConfig() (*rest.Config, error) {
	user := config.getAuthInfo()
	server := config.getServer()

	if err := config.ConfirmUsable(); err != nil {
		return nil, err
	}

	clientConfig := &rest.Config{
		BearerToken:   user.Token,
		Username:      user.Username,
		Password:      user.Password,
		SecretID:      user.SecretID,
		SecretKey:     user.SecretKey,
		Host:          server.Address,
		Timeout:       server.Timeout,
		MaxRetries:    server.MaxRetries,
		RetryInterval: server.RetryInterval,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure:   server.InsecureSkipTLSVerify,
			ServerName: server.TLSServerName,
			CertFile:   user.ClientCertificate,
			KeyFile:    user.ClientKey,
			CertData:   []byte(user.ClientCertificateData),
			KeyData:    []byte(user.ClientKeyData),
			CAFile:     server.CertificateAuthority,
			CAData:     []byte(server.CertificateAuthorityData),
		},
	}

	if u, err := url.ParseRequestURI(clientConfig.Host); err == nil && u.Opaque == "" && len(u.Path) > 1 {
		u.RawQuery = ""
		u.Fragment = ""
		clientConfig.Host = u.String()
	}

	return clientConfig, nil
}
