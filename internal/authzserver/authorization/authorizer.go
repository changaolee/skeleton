package authorization

import (
	"github.com/changaolee/skeleton/internal/pkg/model"
	"github.com/changaolee/skeleton/pkg/log"
	"github.com/ory/ladon"
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
func (a *Authorizer) Authorize(request *ladon.Request) *model.AuthzResponse {
	log.Debugw("authorize request", "request", request)

	if err := a.warden.IsAllowed(request); err != nil {
		return &model.AuthzResponse{
			Denied: true,
			Reason: err.Error(),
		}
	}

	return &model.AuthzResponse{
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

func (a *client) LogRejectedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	// todo: log rejected access request
	log.Infow("TODO: log rejected access request")
}

func (a *client) LogGrantedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	// todo: log granted access request
	log.Infow("TODO: log granted access request")
}
