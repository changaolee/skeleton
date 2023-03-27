// Copyright 2022 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package store

const defaultLimitValue = 20

func defaultLimit(limit int) int {
	if limit == 0 {
		limit = defaultLimitValue
	}
	return limit
}
