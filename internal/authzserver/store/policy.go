// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package store

import "github.com/ory/ladon"

// PolicyStore defines the policy storage interface.
type PolicyStore interface {
	List() (map[string][]*ladon.DefaultPolicy, error)
}
