package auth

type SignupEmailRequest struct {
	Email      string `json:"email" validate:"required,email,email_not_exists"`
	Password   string `json:"password" validate:"required"`
	RePassword string `json:"re_password" validate:"required"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}
