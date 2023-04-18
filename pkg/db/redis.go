package db

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisOptions 定义 Redis 数据库的选项.
type RedisOptions struct {
	Host     string
	Port     int
	Username string
	Password string
	Database int
}

func (o *RedisOptions) Addr() string {
	return fmt.Sprintf("%s:%d", o.Host, o.Port)
}

// NewRedis 使用给定的选项创建一个新的 Redis 数据库实例.
func NewRedis(opts *RedisOptions) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.Addr(),
		Password: opts.Password,
		DB:       opts.Database,
	})
	return rdb, nil
}
