// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package biz

import (
	"context"

	"github.com/changaolee/skeleton/internal/apiserver/store"
	"github.com/changaolee/skeleton/internal/pkg/model/user"
)

type UserBiz interface {
	Create(ctx context.Context, user *user.User) error
	Get(ctx context.Context, username string) (*user.User, error)
}

type userBiz struct {
	s store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func newUsers(b *biz) *userBiz {
	return &userBiz{s: b.s}
}

func (b *userBiz) Create(ctx context.Context, user *user.User) error {
	return b.s.Users().Create(ctx, user)
}

func (b *userBiz) Get(ctx context.Context, username string) (*user.User, error) {
	return b.s.Users().Get(ctx, username)
}
