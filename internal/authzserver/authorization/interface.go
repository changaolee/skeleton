// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package authorization

import "github.com/ory/ladon"

type AuthzInterface interface {
	Create(*ladon.DefaultPolicy) error                    // 创建授权策略
	Update(*ladon.DefaultPolicy) error                    // 更新授权策略
	Delete(id string) error                               // 删除授权策略
	DeleteCollection(idList []string) error               // 批量删除授权策略
	Get(id string) (*ladon.DefaultPolicy, error)          // 获取授权策略
	List(username string) ([]*ladon.DefaultPolicy, error) // 获取指定用户的授权策略列表

	LogRejectedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) // 记录拒绝授权的请求
	LogGrantedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies)  // 记录批准授权的请求
}
