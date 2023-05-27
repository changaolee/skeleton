// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package mysql

import (
	"context"
	"regexp"

	"github.com/changaolee/skeleton/internal/pkg/model/user"
	"gorm.io/gorm"

	"github.com/changaolee/skeleton/internal/apiserver/store"
	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/pkg/errors"
)

type userStore struct {
	ds *datastore
}

var _ store.UserStore = (*userStore)(nil)

func newUsers(ds *datastore) *userStore {
	return &userStore{ds: ds}
}

func (u *userStore) Create(ctx context.Context, user *user.User) error {
	err := u.ds.db.Create(&user).Error
	if err != nil {
		if matched, _ := regexp.MatchString("Duplicate entry '.*' for key 'index_name'", err.Error()); matched {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (u *userStore) Update(ctx context.Context, user *user.User) error {
	err := u.ds.db.Save(user).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (u *userStore) Get(ctx context.Context, username string) (*user.User, error) {
	user := &user.User{}
	err := u.ds.db.Where("name = ? and status = 1", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return user, nil
}
