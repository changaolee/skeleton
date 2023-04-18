package store

import pb "github.com/changaolee/skeleton/pkg/proto/apiserver/v1"

// SecretStore defines the secret storage interface.
type SecretStore interface {
	List() (map[string]*pb.SecretInfo, error)
}
