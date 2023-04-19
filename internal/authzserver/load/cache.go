// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package load

import (
	"sync"

	"github.com/dgraph-io/ristretto"
	"github.com/ory/ladon"

	"github.com/changaolee/skeleton/internal/authzserver/authorization"
	"github.com/changaolee/skeleton/internal/authzserver/store"
	"github.com/changaolee/skeleton/pkg/errors"
	pb "github.com/changaolee/skeleton/pkg/proto/apiserver/v1"
)

// Cache 用于存储 secrets 和 policies.
type Cache struct {
	lock     *sync.RWMutex
	s        store.IStore
	secrets  *ristretto.Cache
	policies *ristretto.Cache
}

// 需要实现的接口.
var _ Loader = &Cache{}
var _ authorization.PolicyGetter = &Cache{}

var (
	ErrSecretNotFound = errors.New("secret not found") // secret 未找到
	ErrPolicyNotFound = errors.New("policy not found") // policy 未找到
)

var (
	onceCache sync.Once
	cacheIns  *Cache
)

// GetCacheInstance 基于指定 store 实例获取 Cache 实例.
func GetCacheInstance(s store.IStore) (*Cache, error) {
	var err error
	if s != nil {
		var (
			secretCache *ristretto.Cache
			policyCache *ristretto.Cache
		)

		onceCache.Do(func() {
			c := &ristretto.Config{
				NumCounters: 1e7,     // number of keys to track frequency of (10M).
				MaxCost:     1 << 30, // maximum cost of cache (1GB).
				BufferItems: 64,      // number of keys per Get buffer.
				Cost:        nil,
			}

			secretCache, err = ristretto.NewCache(c)
			if err != nil {
				return
			}
			policyCache, err = ristretto.NewCache(c)
			if err != nil {
				return
			}

			cacheIns = &Cache{
				s:        s,
				lock:     new(sync.RWMutex),
				secrets:  secretCache,
				policies: policyCache,
			}
		})
	}

	return cacheIns, err
}

// GetSecret 获取指定用户对应的 secret 详情.
func (c *Cache) GetSecret(key string) (*pb.SecretInfo, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.secrets.Get(key)
	if !ok {
		return nil, ErrSecretNotFound
	}

	return value.(*pb.SecretInfo), nil
}

// GetPolicy 获取指定用户的 ladon policy 详情.
func (c *Cache) GetPolicy(key string) ([]*ladon.DefaultPolicy, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.policies.Get(key)
	if !ok {
		return nil, ErrPolicyNotFound
	}

	return value.([]*ladon.DefaultPolicy), nil
}

// Reload 实现 Loader 接口的重载方法，用于重载 secrets 和 policies.
func (c *Cache) Reload() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 重载 secrets
	// todo: 优化为分页加载方式
	secrets, err := c.s.Secrets().List()
	if err != nil {
		return errors.Wrap(err, "list secrets failed")
	}
	c.secrets.Clear()
	for key, val := range secrets {
		c.secrets.Set(key, val, 1)
	}

	// 重载 policies
	// todo: 优化为分页加载方式
	policies, err := c.s.Policies().List()
	if err != nil {
		return errors.Wrap(err, "list policies failed")
	}
	c.policies.Clear()
	for key, val := range policies {
		c.policies.Set(key, val, 1)
	}

	return nil
}
