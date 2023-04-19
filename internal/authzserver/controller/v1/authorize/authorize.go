// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package authorize

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/ladon"

	"github.com/changaolee/skeleton/internal/authzserver/authorization"
	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/pkg/errors"
)

type AuthzController struct {
	getter authorization.PolicyGetter
}

// NewAuthzController 创建一个 authz controller.
func NewAuthzController(getter authorization.PolicyGetter) *AuthzController {
	return &AuthzController{getter: getter}
}

// Authorize 返回一个资源是否被允许访问.
func (a *AuthzController) Authorize(c *gin.Context) {
	var r ladon.Request
	if err := c.ShouldBind(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	auth := authorization.NewAuthorizer(a.getter)
	if r.Context == nil {
		r.Context = ladon.Context{}
	}

	r.Context["username"] = c.GetString("username")
	rsp := auth.Authorize(&r)

	core.WriteResponse(c, nil, rsp)
}
