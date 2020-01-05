package types

import (
	"errors"
)

var (
	ErrUsernameTaken         error = errors.New("Username is taken")
	ErrEmailTaken            error = errors.New("EMail is taken")
	ErrUserRegister          error = errors.New("Failed to Register User")
	ErrPasswordForgot        error = errors.New("Failed to complete password reset process")
	ErrPasswordResetValidate error = errors.New("Failed to validate a password reset token")
	ErrPasswordReset         error = errors.New("Failed to complete the password reset")
)
