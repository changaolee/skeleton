package client

import (
	"fmt"
	"os"
	"reflect"

	"github.com/changaolee/skeleton/pkg/errors"
	"go.uber.org/multierr"
)

var (
	// ErrEmptyConfig defines no configuration has been provided error.
	ErrEmptyConfig = errors.New("no configuration has been provided, try setting SKT_SERVER_ADDRESS environment variable")

	// ErrEmptyServer defines a no server defined error.
	ErrEmptyServer = errors.New("server has no server defined")
)

type errConfigurationInvalid []error

func newErrConfigurationInvalid(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	default:
		return errConfigurationInvalid(errs)
	}
}

func (e errConfigurationInvalid) Error() string {
	var aggerr error
	for _, err := range e {
		aggerr = multierr.Append(aggerr, err)
	}
	return fmt.Sprintf("invalid configuration: %v", aggerr.Error())
}

// validateServerInfo looks for conflicts and errors in the server info.
func validateServerInfo(serverInfo Server) []error {
	validationErrors := make([]error, 0)

	emptyServer := &Server{}
	if reflect.DeepEqual(*emptyServer, serverInfo) {
		return []error{ErrEmptyServer}
	}

	/*
		if len(serverInfo.Address) == 0 {
			validationErrors = append(validationErrors, fmt.Errorf("no server found"))
		}
	*/
	// Make sure CA data and CA file aren't both specified
	if len(serverInfo.CertificateAuthority) != 0 && len(serverInfo.CertificateAuthorityData) != 0 {
		validationErrors = append(
			validationErrors,
			fmt.Errorf(
				"certificate-authority-data and certificate-authority are both specified. certificate-authority-data will override",
			),
		)
	}

	if len(serverInfo.CertificateAuthority) != 0 {
		clientCertCA, err := os.Open(serverInfo.CertificateAuthority)
		if err != nil {
			validationErrors = append(validationErrors,
				fmt.Errorf("unable to read certificate-authority %v due to %w", serverInfo.CertificateAuthority, err))
		} else {
			defer clientCertCA.Close()
		}
	}

	return validationErrors
}

func getAuthMethods(authInfo AuthInfo) []string {
	methods := make([]string, 0, 3)
	if len(authInfo.Token) != 0 {
		methods = append(methods, "token")
	}

	if len(authInfo.Username) != 0 || len(authInfo.Password) != 0 {
		methods = append(methods, "basicAuth")
	}

	if len(authInfo.SecretID) != 0 || len(authInfo.SecretKey) != 0 {
		methods = append(methods, "secretAuth")
	}
	return methods
}

// validateAuthInfo looks for conflicts and errors in the auth info.
func validateAuthInfo(authInfo AuthInfo) []error {
	validationErrors := make([]error, 0)

	// authPath also provides information for the client to identify the server,
	// so allow multiple auth methods in that case
	methods := getAuthMethods(authInfo)
	if len(methods) > 1 {
		validationErrors = append(validationErrors,
			fmt.Errorf("more than one authentication method found; found %v, only one is allowed", methods))
	}

	if len(authInfo.ClientCertificate) == 0 || len(authInfo.ClientCertificateData) == 0 {
		return validationErrors
	}

	// Make sure cert data and file aren't both specified
	if len(authInfo.ClientCertificate) != 0 && len(authInfo.ClientCertificateData) != 0 {
		validationErrors = append(validationErrors,
			fmt.Errorf("client-cert-data and client-cert are both specified. client-cert-data will override"))
	}
	// Make sure key data and file aren't both specified
	if len(authInfo.ClientKey) != 0 && len(authInfo.ClientKeyData) != 0 {
		validationErrors = append(validationErrors,
			fmt.Errorf("client-key-data and client-key are both specified; client-key-data will override"))
	}
	// Make sure a key is specified
	if len(authInfo.ClientKey) == 0 && len(authInfo.ClientKeyData) == 0 {
		validationErrors = append(validationErrors,
			fmt.Errorf("client-key-data or client-key must be specified to use the clientCert authentication method"))
	}

	if len(authInfo.ClientCertificate) != 0 {
		clientCertFile, err := os.Open(authInfo.ClientCertificate)
		if err != nil {
			validationErrors = append(validationErrors,
				fmt.Errorf("unable to read client-cert %v due to %w", authInfo.ClientCertificate, err))
		} else {
			defer clientCertFile.Close()
		}
	}

	if len(authInfo.ClientKey) != 0 {
		clientKeyFile, err := os.Open(authInfo.ClientKey)
		if err != nil {
			validationErrors = append(validationErrors,
				fmt.Errorf("unable to read client-key %v due to %w", authInfo.ClientKey, err))
		} else {
			defer clientKeyFile.Close()
		}
	}

	return validationErrors
}
