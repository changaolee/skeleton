package apiserver

import (
	"github.com/changaolee/skeleton/internal/authzserver/store"
	pb "github.com/changaolee/skeleton/pkg/proto/apiserver/v1"
)

type datastore struct {
	cli pb.CacheClient
}

var _ store.IStore = (*datastore)(nil)

func (d *datastore) Policies() store.PolicyStore {
	//TODO implement me
	panic("implement me")
}

func (d *datastore) Secrets() store.SecretStore {
	//TODO implement me
	panic("implement me")
}

func GetAPIServerInstance(address, clientCA string) store.IStore {
	// todo: apiserver instance
}
