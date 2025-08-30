package dto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	MobileNumber string `json:"email" validate:"required,min=11,max=11"`
}

type LoginRequest struct {
	MobileNumber string `json:"email" validate:"required,min=11,max=11"`
	OTP          int    `json:"otp" validate:"required,min=6,max=6"`
}

func ValidateModel(model any) error {
	validate := validator.New()
	if err := validate.Struct(model); err != nil {
		var sb strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			sb.WriteString(fmt.Sprintf("Field '%s' failed on '%s'\n", err.Field(), err.Tag()))
		}
		return errors.New(sb.String())
	}
	return nil
}
