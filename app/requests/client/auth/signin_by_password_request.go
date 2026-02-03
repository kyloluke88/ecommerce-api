package auth

type SigninByPasswordRequest struct {
	LoginId  string `json:"login_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}
