// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package v1

import (
	"time"
)

// ObjectMeta 是所有持久化资源必须包含的字段.
type ObjectMeta struct {
	ID           uint64    `json:"id,omitempty"         gorm:"primary_key;AUTO_INCREMENT;column:id"`
	InstanceID   string    `json:"instanceID,omitempty" gorm:"unique;column:instanceID;type:varchar(32);not null"`
	Extend       Extend    `json:"extend,omitempty"     gorm:"-"                                                  validate:"omitempty"`
	ExtendShadow string    `json:"-"                    gorm:"column:extendShadow"                                validate:"omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"  gorm:"column:createdAt"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"  gorm:"column:updatedAt"`
}
