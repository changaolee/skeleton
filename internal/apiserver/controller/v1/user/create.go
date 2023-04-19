// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package user

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/pkg/errors"

	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/internal/pkg/model"
	"github.com/changaolee/skeleton/pkg/log"
)

// Create 创建一个新用户.
func (u *UserController) Create(c *gin.Context) {
	log.C(c).Infow("Create user function called")

	var r model.User
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, err.Error()), nil)
		return
	}

	r.Status = 1
	r.LoginAt = time.Now()

	if err := u.b.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, r)
}
