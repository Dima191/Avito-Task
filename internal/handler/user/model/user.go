package userhandlermodel

import (
	"github.com/go-playground/validator/v10"
	"slices"
)

type User struct {
	Role     string `json:"role" validate:"required,role"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

var (
	PossibleRoles = []string{"client", "moderator"}
)

func PasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	switch {
	case len(password) < 6:
		return false
	default:
		return true
	}
}

func RoleValidation(fl validator.FieldLevel) bool {
	return slices.Contains(PossibleRoles, fl.Field().String())
}
