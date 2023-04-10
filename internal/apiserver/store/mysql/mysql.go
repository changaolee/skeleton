package mysql

import (
	"fmt"
	"sync"

	"github.com/changaolee/skeleton/internal/apiserver/store"
	genoptions "github.com/changaolee/skeleton/internal/pkg/options"
	"github.com/changaolee/skeleton/pkg/db"
	"gorm.io/gorm"
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

	var err error
	var dbIns *gorm.DB

	if mysqlIns == nil {
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
			dbIns, err = db.NewMySQL(options)
			mysqlIns = &datastore{db: dbIns}
		})
	}
	if mysqlIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql instance, mysqlIns: %+v, error: %w", mysqlIns, err)
	}
	return mysqlIns, nil
}
