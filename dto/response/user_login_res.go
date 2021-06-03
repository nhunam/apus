package response

type UserLoginRes struct {
	UserId     uint64
	CompanyId  uint64
	RememberMe *bool
}
