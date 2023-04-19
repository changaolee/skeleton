// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package options

import (
	"github.com/spf13/pflag"
)

// RedisOptions 定义了 Redis 选项.
type RedisOptions struct {
	Host     string `json:"host"     mapstructure:"host"`
	Port     int    `json:"port"     mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Database int    `json:"database" mapstructure:"database"`
}

// NewRedisOptions 创建一个默认值的 Redis 选项实例.
func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Host:     "127.0.0.1",
		Port:     6379,
		Username: "",
		Password: "",
		Database: 0,
	}
}

// Validate 验证 Redis 选项.
func (o *RedisOptions) Validate() []error {
	var errs []error

	return errs
}

// AddFlags 向指定 FlagSet 中添加 Redis 选项相关标志.
func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "redis.host", o.Host, "Hostname of your Redis server.")
	fs.IntVar(&o.Port, "redis.port", o.Port, "The port the Redis server is listening on.")
	fs.StringVar(&o.Username, "redis.username", o.Username, "Username for access to redis service.")
	fs.StringVar(&o.Password, "redis.password", o.Password, "Optional auth password for Redis db.")

	fs.IntVar(&o.Database, "redis.database", o.Database, ""+
		"By default, the database is 0. Setting the database is not supported with redis cluster. "+
		"As such, if you have --redis.enable-cluster=true, then this value should be omitted or explicitly set to 0.")
}
