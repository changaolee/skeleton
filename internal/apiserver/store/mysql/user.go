package mysql

import (
	"context"

	"github.com/changaolee/skeleton/internal/apiserver/store"
	"github.com/changaolee/skeleton/internal/pkg/model"
)

type users struct {
	ds *datastore
}

var _ store.UserStore = (*users)(nil)

func newUsers(ds *datastore) *users {
	return &users{ds: ds}
}

func (u *users) Create(ctx context.Context, user *model.User) error {
	return u.ds.db.Create(&user).Error
}
