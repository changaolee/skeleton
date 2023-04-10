package v1

import (
	"time"
)

// ObjectMeta 是所有持久化资源必须包含的字段.
type ObjectMeta struct {
	ID           uint64    `json:"id,omitempty"  gorm:"primary_key;AUTO_INCREMENT;column:id"`
	InstanceID   string    `json:"instanceID,omitempty" gorm:"unique;column:instanceID;type:varchar(32);not null"`
	Extend       Extend    `json:"extend,omitempty" gorm:"-" validate:"omitempty"`
	ExtendShadow string    `json:"-" gorm:"column:extendShadow" validate:"omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty" gorm:"column:createdAt"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty" gorm:"column:updatedAt"`
}
