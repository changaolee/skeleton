// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package biz

import "github.com/changaolee/skeleton/internal/apiserver/store"

// IBiz 定义了 Biz 层接口.
type IBiz interface {
	Users() UserBiz
}

type biz struct {
	s store.IStore
}

var _ IBiz = (*biz)(nil)

// New 创建一个.
func New(s store.IStore) *biz {
	return &biz{s: s}
}

func (b *biz) Users() UserBiz {
	return newUsers(b)
}
