// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package mysql

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"github.com/changaolee/skeleton/internal/apiserver/store"
	genoptions "github.com/changaolee/skeleton/internal/pkg/options"
	"github.com/changaolee/skeleton/pkg/db"
)

type datastore struct {
	db *gorm.DB
}

var _ store.IStore = (*datastore)(nil)

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Close() error {
	conn, err := ds.db.DB()
	if err != nil {
		return fmt.Errorf("get gorm db instance failed")
	}
	return conn.Close()
}

var (
	mysqlIns store.IStore
	once     sync.Once
)

// GetMySQLInstance 获取 MySQL 实例.
func GetMySQLInstance(opts *genoptions.MySQLOptions) (store.IStore, error) {
	if opts == nil && mysqlIns == nil {
		return nil, fmt.Errorf("failed to get mysql instance")
	}

	var (
		err error
		ins *gorm.DB
	)

	once.Do(func() {
		options := &db.MySQLOptions{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
		}
		ins, err = db.NewMySQL(options)
		mysqlIns = &datastore{db: ins}
	})

	if mysqlIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql instance, mysqlIns: %+v, error: %w", mysqlIns, err)
	}
	return mysqlIns, nil
}
