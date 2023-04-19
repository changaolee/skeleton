package apiserver

import (
	"context"

	"github.com/AlekSi/pointer"
	"github.com/avast/retry-go"
	"github.com/changaolee/skeleton/internal/authzserver/store"
	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/changaolee/skeleton/pkg/log"
	pb "github.com/changaolee/skeleton/pkg/proto/apiserver/v1"
)

type secretStore struct {
	ds *datastore
}

var _ store.SecretStore = (*secretStore)(nil)

func newSecrets(ds *datastore) *secretStore {
	return &secretStore{ds: ds}
}

// List 返回所有的授权密钥.
func (s *secretStore) List() (map[string]*pb.SecretInfo, error) {
	secrets := make(map[string]*pb.SecretInfo)

	log.Infow("Loading secrets")

	req := &pb.ListSecretsRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	var resp *pb.ListSecretsResponse
	err := retry.Do(
		func() error {
			var listErr error
			resp, listErr = s.ds.cli.ListSecrets(context.Background(), req)
			if listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		return nil, errors.Wrap(err, "list secrets failed")
	}

	log.Infof("Secrets found (%d total):", len(resp.Items))

	for _, v := range resp.Items {
		log.Infof(" - %s:%s", v.Username, v.SecretId)
		secrets[v.SecretId] = v
	}

	return secrets, nil
}
