package model

import (
	"time"
)

type Otp struct {
	ID        uint64     `gorm:"column:id;primaryKey"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	Type      string     `gorm:"column:type;not null"`
	Otp       string     `gorm:"column:otp;not null"`
	UserId    *int64     `gorm:"column:user_id"`
	ValidTo   *time.Time `gorm:"column:valid_to"`
	Valid     *bool      `gorm:"column:is_valid"`
}

func (Otp) TableName() string {
	return "otps"
}
