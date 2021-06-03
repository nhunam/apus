package response

import (
	"time"
)

type RefreshTokenRes struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiredTime  time.Duration `json:"expired_time"`
}
