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

type userStore struct {
	ds *datastore
}

var _ store.UserStore = (*userStore)(nil)

func newUsers(ds *datastore) *userStore {
	return &userStore{ds: ds}
}

func (u *userStore) Create(ctx context.Context, user *model.User) error {
	return u.ds.db.Create(&user).Error
}