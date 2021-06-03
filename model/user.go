package model

import (
	"time"
)

type User struct {
	ID               uint64           `gorm:"column:id;primaryKey"`
	CreatedAt        *time.Time       `gorm:"column:created_at;autoCreateTime"`
	CreatedBy        *uint64          `gorm:"column:created_by"`
	UpdatedAt        *time.Time       `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy        *uint64          `gorm:"column:updated_by"`
	DeletedAt        *time.Time       `gorm:"column:deleted_at"`
	FirstName        *string          `gorm:"column:first_name"`
	MiddleName       *string          `gorm:"column:middle_name"`
	LastName         *string          `gorm:"column:last_name"`
	Address          *string          `gorm:"column:address"`
	Gender           *string          `gorm:"column:gender"`
	Birthday         *time.Time       `gorm:"column:birthday,type:date"`
	IdNumber         *string          `gorm:"column:id_number"`
	IdCardFront      *string          `gorm:"column:id_card_front"`
	IdCardBack       *string          `gorm:"column:id_card_back"`
	DfCompanyId      *int64           `gorm:"column:df_company_id"`
	VerifyStatus     string           `gorm:"column:verify_status;not null"`
	ActivateStatus   string           `gorm:"column:activate_status;not null"`
	DfLanguageCode   *string          `gorm:"column:df_language_code"`
	DfCountryCode    *string          `gorm:"column:df_country_code"`
	DfCurrencyCode   *string          `gorm:"column:df_currency_code"`
	DfTimezoneOffset *float32         `gorm:"column:df_timezone_offset"`
	DfTimezoneAbbr   *string          `gorm:"column:df_timezone_abbr"`
	ThemeName        *string          `gorm:"column:theme_name"`
	ThemeWriteSys    *string          `gorm:"column:theme_write_sys"`
	Attributes       []UserAttribute  `gorm:"foreignKey:UserId;constraint:OnCreate:CASCADE,OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Credentials      []UserCredential `gorm:"foreignKey:UserId;constraint:OnCreate:CASCADE,OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (User) TableName() string {
	return "users"
}
