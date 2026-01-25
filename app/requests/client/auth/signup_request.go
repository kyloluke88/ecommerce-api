package auth

import (
	"api/app/models/user"

	"github.com/go-playground/validator/v10"
)

type SignupRequest struct {
	Email string `json:"email" validate:"required,email,email_not_exists"`
	Phone string `json:"phone" validate:"required"`
}

func RegisterValidations(v *validator.Validate) {
	_ = v.RegisterValidation("email_not_exists", emailNotExists)
}


func emailNotExists(fl validator.FieldLevel) bool {
	return !user.IsEmailExist(fl.Field().String())
}
