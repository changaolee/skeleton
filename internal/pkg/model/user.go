// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	metav1 "github.com/changaolee/skeleton/pkg/meta/v1"
	"github.com/changaolee/skeleton/pkg/util/idutil"

	"github.com/changaolee/skeleton/pkg/auth"
)

// User 是数据库中 user 记录 struct 格式的映射.
type User struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            int       `json:"status"             gorm:"column:status"   validate:"omitempty"`
	Nickname          string    `json:"nickname"           gorm:"column:nickname" validate:"required,min=1,max=30"`
	Password          string    `json:"password,omitempty" gorm:"column:password" validate:"required"`
	Email             string    `json:"email"              gorm:"column:email"    validate:"required,email,min=1,max=100"`
	Phone             string    `json:"phone"              gorm:"column:phone"    validate:"omitempty"`
	IsAdmin           int       `json:"isAdmin,omitempty"  gorm:"column:isAdmin"  validate:"omitempty"`
	TotalPolicy       int64     `json:"totalPolicy"        gorm:"-"               validate:"omitempty"`
	LoginAt           time.Time `json:"loginAt,omitempty"  gorm:"column:loginAt"`
}

// TableName 用来指定映射的 MySQL 表名.
func (u *User) TableName() string {
	return "user"
}

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}
	return nil
}

// AfterCreate 在创建数据库记录之后更新资源 ID.
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "user-")

	return tx.Save(u).Error
}

// Compare 验证用户密码是否正确.
func (u *User) Compare(pwd string) error {
	if err := auth.Compare(u.Password, pwd); err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}
	return nil
}
