package request

type UserLoginDto struct {
	Mobile     *string `json:"mobile" binding:"omitempty,numeric,max=20"`
	Email      *string `json:"email" binding:"omitempty,email"`
	Username   *string `json:"username" biding:"omitempty,min=5,max=12"`
	Password   string  `json:"password" binding:"required,min=6,max=12"`
	RememberMe *bool   `json:"remember_me"`
}
