// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package apiserver

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/changaolee/skeleton/internal/authzserver/store"
	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/changaolee/skeleton/pkg/log"
	pb "github.com/changaolee/skeleton/pkg/proto/apiserver/v1"
)

type datastore struct {
	cli pb.CacheClient
}

var _ store.IStore = (*datastore)(nil)

func (d *datastore) Policies() store.PolicyStore {
	return newPolicies(d)
}

func (d *datastore) Secrets() store.SecretStore {
	return newSecrets(d)
}

var (
	cacheIns store.IStore
	once     sync.Once
)

// GetAPIServerCacheClientInstance 获取 APIServer CacheClient 实例.
func GetAPIServerCacheClientInstance(address, clientCA string) (store.IStore, error) {
	var (
		err   error
		conn  *grpc.ClientConn
		creds credentials.TransportCredentials
	)

	once.Do(func() {
		creds, err = credentials.NewClientTLSFromFile(clientCA, "")
		if err != nil {
			err = errors.Wrap(err, "credentials.NewClientTLSFromFile error")
			return
		}

		conn, err = grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
		if err != nil {
			err = errors.Wrap(err, "Connect to grpc server failed")
			return
		}

		cacheIns = &datastore{cli: pb.NewCacheClient(conn)}
		log.Infof("Connected to grpc server, address: %s", address)
	})

	if cacheIns == nil || err != nil {
		return nil, errors.Wrapf(err, "failed to get cache client instance, cacheIns: %+v", cacheIns)
	}
	return cacheIns, err
}
