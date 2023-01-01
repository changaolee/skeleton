// Copyright 2022 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package user

import (
	"context"
	"regexp"

	"github.com/jinzhu/copier"

	"github.com/changaolee/skeleton/internal/pkg/errno"
	"github.com/changaolee/skeleton/internal/pkg/model"
	"github.com/changaolee/skeleton/internal/skeleton/store"
	v1 "github.com/changaolee/skeleton/pkg/api/skeleton/v1"
)

// UserBiz 定义了 user 模块在 biz 层所实现的方法
type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
}

// UserBiz 接口的实现
type userBiz struct {
	ds store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口
var _ UserBiz = (*userBiz)(nil)

// New 创建一个实现了 UserBiz 接口的实例
func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

// Create 是 UserBiz 接口中 `Create` 方法的实现
func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}
	return nil
}
