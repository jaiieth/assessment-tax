package helper

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return errors.New("invalid request")
	}
	return nil
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New(validator.WithRequiredStructEnabled())}
}
