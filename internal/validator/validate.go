package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
)

const (
	EmailTag            = "email"
	PasswordTag         = "password"
	RoleTag             = "role"
	ModerationStatusTag = "moderation_status"
)

var (
	ErrInvalidEmail            = errors.New("invalid email")
	ErrInvalidPassword         = errors.New("password validation failed. length should be at least 6 characters")
	ErrInvalidRole             = errors.New("invalid role. possible roles: client, moderator")
	ErrInvalidModerationStatus = errors.New("invalid moderation status. possible status: created, approved, declined, on moderation")
)

type Validate struct {
	validate *validator.Validate
}

func (v *Validate) Validate(model any) (resErr error) {
	err := v.validate.Struct(model)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationErr := range validationErrors {
				switch validationErr.Tag() {
				case EmailTag:
					resErr = multierror.Append(resErr, ErrInvalidEmail)
				case PasswordTag:
					resErr = multierror.Append(resErr, ErrInvalidPassword)
				case RoleTag:
					resErr = multierror.Append(resErr, ErrInvalidRole)
				case ModerationStatusTag:
					resErr = multierror.Append(resErr, ErrInvalidModerationStatus)
				default:
					resErr = multierror.Append(resErr, err)
				}
			}
		}
		return resErr
	}

	return nil
}

func (v *Validate) RegisterTag(tag string, fn validator.Func) error {
	if err := v.validate.RegisterValidation(tag, fn); err != nil {
		return err
	}
	return nil
}

func New() *Validate {
	return &Validate{
		validate: validator.New(),
	}
}
