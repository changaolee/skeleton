// Copyright 2022 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package user

import (
	"github.com/changaolee/skeleton/internal/skeleton/biz"
	"github.com/changaolee/skeleton/internal/skeleton/store"
	pb "github.com/changaolee/skeleton/pkg/proto/skeleton/v1"
)

// UserController 是 user 模块在 Controller 层的实现，用来处理用户模块的请求.
type UserController struct {
	b biz.IBiz
	pb.UnimplementedSkeletonServer
}

// New 创建一个 user controller.
func New(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}
