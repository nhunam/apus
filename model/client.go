package model

import (
	"time"
)

type Client struct {
	ID                   uint64     `gorm:"column:id;primaryKey"`
	CreatedAt            *time.Time `gorm:"column:created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at"`
	DeletedAt            *time.Time `gorm:"column:deleted_at"`
	ClientId             string     `gorm:"column:client_id;not null"`
	ClientSecret         string     `gorm:"column:client_secret;not null"`
	AccessTokenValidity  *int64     `gorm:"column:access_token_validity"`
	RefreshTokenValidity *int64     `gorm:"column:refresh_token_validity"`
	Scope                *string    `gorm:"column:scope;not null"`
}

func (Client) TableName() string {
	return "clients"
}
