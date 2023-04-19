// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package store

// IStore 定义了 Store 层接口.
type IStore interface {
	Policies() PolicyStore
	Secrets() SecretStore
}

var ins IStore

// Store 获取 Store 实例.
func Store() IStore {
	return ins
}

// SetStore 设置 Store 实例.
func SetStore(storeIns IStore) {
	ins = storeIns
}
