// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

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
