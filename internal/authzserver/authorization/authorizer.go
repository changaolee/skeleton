// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package authorization

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/changaolee/skeleton/internal/pkg/response"
	"github.com/ory/ladon"

	"github.com/changaolee/skeleton/internal/authzserver/analytics"

	"github.com/changaolee/skeleton/pkg/log"
)

// PolicyGetter 定义获取指定用户授权策略的接口.
type PolicyGetter interface {
	GetPolicy(key string) ([]*ladon.DefaultPolicy, error)
}

// Authorizer 实现了授权审核接口.
type Authorizer struct {
	warden ladon.Warden
}

// NewAuthorizer 创建一个 Authorizer 实例.
func NewAuthorizer(getter PolicyGetter) *Authorizer {
	cli := &client{getter: getter}
	return &Authorizer{
		warden: &ladon.Ladon{
			Manager:     NewPolicyManager(cli),
			AuditLogger: NewAuditLogger(cli),
		},
	}
}

// Authorize 确定访问权限.
func (a *Authorizer) Authorize(request *ladon.Request) *response.AuthzResponse {
	log.Debugw("authorize request", "request", request)

	if err := a.warden.IsAllowed(request); err != nil {
		return &response.AuthzResponse{
			Denied: true,
			Reason: err.Error(),
		}
	}

	return &response.AuthzResponse{
		Allowed: true,
	}
}

type client struct {
	getter PolicyGetter
}

var _ AuthzInterface = (*client)(nil)

func (a *client) Create(policy *ladon.DefaultPolicy) error {
	return nil
}

func (a *client) Update(policy *ladon.DefaultPolicy) error {
	return nil
}

func (a *client) Delete(id string) error {
	return nil
}

func (a *client) DeleteCollection(idList []string) error {
	return nil
}

func (a *client) Get(id string) (*ladon.DefaultPolicy, error) {
	return &ladon.DefaultPolicy{}, nil
}

func (a *client) List(username string) ([]*ladon.DefaultPolicy, error) {
	return a.getter.GetPolicy(username)
}

func (a *client) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	var conclusion string
	if len(d) > 1 {
		allowed := joinPoliciesNames(d[0 : len(d)-1])
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("policies %s allow access, but policy %s forcefully denied it", allowed, denied)
	} else if len(d) == 1 {
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("policy %s forcefully denied the access", denied)
	} else {
		conclusion = "no policy allowed access"
	}
	rstring, pstring, dstring := convertToString(r, p, d)
	record := analytics.Record{
		TimeStamp:  time.Now().Unix(),
		Username:   r.Context["username"].(string),
		Effect:     ladon.DenyAccess,
		Conclusion: conclusion,
		Request:    rstring,
		Policies:   pstring,
		Deciders:   dstring,
	}

	// todo: local cache + redis
	log.Infof("Log rejected access request: %+v", record)
}

func (a *client) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	conclusion := fmt.Sprintf("policies %s allow access", joinPoliciesNames(d))
	rstring, pstring, dstring := convertToString(r, p, d)
	record := analytics.Record{
		TimeStamp:  time.Now().Unix(),
		Username:   r.Context["username"].(string),
		Effect:     ladon.AllowAccess,
		Conclusion: conclusion,
		Request:    rstring,
		Policies:   pstring,
		Deciders:   dstring,
	}

	// todo: local cache + redis
	log.Infof("Log granted access request: %+v", record)
}

func joinPoliciesNames(policies ladon.Policies) string {
	names := make([]string, 0, len(policies))
	for _, policy := range policies {
		names = append(names, policy.GetID())
	}

	return strings.Join(names, ", ")
}

func convertToString(r *ladon.Request, p ladon.Policies, d ladon.Policies) (string, string, string) {
	rbytes, _ := json.Marshal(r)
	pbytes, _ := json.Marshal(p)
	dbytes, _ := json.Marshal(d)

	return string(rbytes), string(pbytes), string(dbytes)
}
