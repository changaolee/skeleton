// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package user

import (
	"github.com/changaolee/skeleton/internal/apiserver/biz"
	"github.com/changaolee/skeleton/internal/apiserver/store"
)

type UserController struct {
	b biz.IBiz
}

// New 创建一个 user controller.
func New(s store.IStore) *UserController {
	return &UserController{b: biz.New(s)}
}
