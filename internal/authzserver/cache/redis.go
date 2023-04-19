// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"

	genoptions "github.com/changaolee/skeleton/internal/pkg/options"
	"github.com/changaolee/skeleton/pkg/db"
	"github.com/changaolee/skeleton/pkg/log"
)

type RedisCache struct {
	rd *redis.Client
}

var (
	redisIns *RedisCache
	once     sync.Once
)

// GetRedisInstance 获取 Redis 实例.
func GetRedisInstance(opts *genoptions.RedisOptions) (*RedisCache, error) {
	if opts == nil && redisIns == nil {
		return nil, fmt.Errorf("failed to get redis instance")
	}

	var err error
	var ins *redis.Client

	if redisIns == nil {
		once.Do(func() {
			options := &db.RedisOptions{
				Host:     opts.Host,
				Port:     opts.Port,
				Username: opts.Username,
				Password: opts.Password,
				Database: opts.Database,
			}
			ins, err = db.NewRedis(options)
			redisIns = &RedisCache{rd: ins}
		})
	}
	if redisIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get redis instance, redisIns: %+v, error: %w", redisIns, err)
	}
	return redisIns, nil
}

func (r *RedisCache) StartPubSubHandler(ctx context.Context, channel string, callback func(interface{})) error {
	subscriber := r.rd.Subscribe(ctx, channel)
	defer subscriber.Close()

	if _, err := subscriber.Receive(ctx); err != nil {
		log.Errorf("Error while receiving pubsub message: %s", err.Error())
		return err
	}

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			log.Errorf("Error while receiving pubsub message: %s", err.Error())
		}
		callback(msg)
	}
}
