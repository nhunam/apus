package model

import (
	"time"
)

type I18nMessageId struct {
	Code         string `gorm:"column:code;not null"`
	LanguageCode string `gorm:"column:language_code;not null"`
	CountryCode  string `gorm:"column:country_code;not null"`
}

type I18nMessage struct {
	ID        I18nMessageId `gorm:"embedded"`
	CreatedAt *time.Time    `gorm:"column:created_at"`
	UpdatedAt *time.Time    `gorm:"column:updated_at"`
	DeletedAt *time.Time    `gorm:"column:deleted_at"`
	Message   string        `gorm:"column:message"`
}

func (I18nMessage) TableName() string {
	return "i18_messages"
}
