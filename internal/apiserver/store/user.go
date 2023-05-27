// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package store

import (
	"context"

	"github.com/changaolee/skeleton/internal/pkg/model/user"
)

type UserStore interface {
	Create(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
	Get(ctx context.Context, username string) (*user.User, error)
}
