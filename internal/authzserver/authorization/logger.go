package authorization

import (
	"github.com/changaolee/skeleton/pkg/log"
	"github.com/ory/ladon"
)

// AuditLogger 输出并缓存批准或拒绝授权的日志.
type AuditLogger struct {
	client AuthzInterface
}

// NewAuditLogger 创建一个 AuditLogger 实例.
func NewAuditLogger(client AuthzInterface) ladon.AuditLogger {
	return &AuditLogger{client: client}
}

func (l *AuditLogger) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	l.client.LogRejectedAccessRequest(r, p, d)
	log.Debugw("Subject access review rejected", "request", r, "deciders", d)
}

func (l *AuditLogger) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	l.client.LogGrantedAccessRequest(r, p, d)
	log.Debugw("Subject access review granted", "request", r, "deciders", d)
}
