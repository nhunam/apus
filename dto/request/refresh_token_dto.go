package request

type RefreshTokenDto struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	CompanyId    string `json:"company_id" binding:"required"`
}
