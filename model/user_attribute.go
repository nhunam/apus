package model

import (
	"time"
)

type UserAttribute struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	CreatedAt   *time.Time `gorm:"column:created_at;autoCreateTime"`
	CreatedBy   *uint64    `gorm:"column:created_by"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy   *uint64    `gorm:"column:updated_by"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	UserId      int64      `gorm:"column:user_id;not null"`
	Type        string     `gorm:"column:type;not null"`
	Code        string     `gorm:"column:code;not null"`
	Name        *string    `gorm:"column:name;not null"`
	ValueString *string    `gorm:"column:value_string"`
}

func (UserAttribute) TableName() string {
	return "user_attrs"
}
