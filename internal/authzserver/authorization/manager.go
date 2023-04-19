// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package authorization

import (
	"github.com/ory/ladon"

	"github.com/changaolee/skeleton/pkg/errors"
)

// PolicyManager 是一个基于 MySQL 实现的授权策略持久化 Manager.
type PolicyManager struct {
	client AuthzInterface
}

// NewPolicyManager 创建一个 PolicyManager 实例.
func NewPolicyManager(client AuthzInterface) ladon.Manager {
	return &PolicyManager{client: client}
}

func (m *PolicyManager) Create(policy ladon.Policy) error {
	return nil
}

func (m *PolicyManager) Update(policy ladon.Policy) error {
	return nil
}

func (m *PolicyManager) Get(id string) (ladon.Policy, error) {
	return &ladon.DefaultPolicy{}, nil
}

func (m *PolicyManager) Delete(id string) error {
	return nil
}

func (m *PolicyManager) GetAll(limit, offset int64) (ladon.Policies, error) {
	return nil, nil
}

func (m *PolicyManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	username := ""

	if user, ok := r.Context["username"].(string); ok {
		username = user
	}

	policies, err := m.client.List(username)
	if err != nil {
		return nil, errors.Wrap(err, "list policies failed")
	}

	ret := make([]ladon.Policy, 0, len(policies))
	for _, policy := range policies {
		ret = append(ret, policy)
	}

	return ret, nil
}

func (m *PolicyManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	return nil, nil
}

func (m *PolicyManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	return nil, nil
}
