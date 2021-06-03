package model

import (
	"gorm.io/gorm"
	"time"
)

type UserCredential struct {
	ID         uint64          `gorm:"column:id;primaryKey"`
	CreatedAt  *time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  *time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  *gorm.DeletedAt `gorm:"column:deleted_at"`
	UserId     int64           `gorm:"column:user_id;not null"`
	Type       string          `gorm:"column:type;not null"`
	Credential string          `gorm:"column:credential;not null"`
	Password   string          `gorm:"column:password;not null"`
	Active     *bool           `gorm:"column:is_active;type:boolean;"`
}

func (UserCredential) TableName() string {
	return "user_credentials"
}
