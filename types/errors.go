package types

import (
	"errors"
)

var (
	ErrUsernameTaken                  error = errors.New("Username is taken")
	ErrEmailTaken                     error = errors.New("EMail is taken")
	ErrUserRegister                   error = errors.New("Failed to Register User")
	ErrPasswordForgot                 error = errors.New("Failed to complete password reset process")
	ErrPasswordResetValidate          error = errors.New("Failed to validate a password reset token")
	ErrPasswordResetValidateServerErr error = errors.New("Failed to validate a password reset token: Server Error")
	ErrPasswordReset                  error = errors.New("Failed to complete the password reset")
	ErrPasswordResetServerErr         error = errors.New("Failed to complete the password reset: Server Error")

	ErrGetWantedItemsDB   error = errors.New("Failed to Get Wanted Items from DB")
	ErrGetAllItemsDB      error = errors.New("Failed to Get All Items from DB")
	ErrGetReservedItemsDB error = errors.New("Failed to Get Reserved Items from DB")
	ErrAddItemErr         error = errors.New("Failed to Add Item")
	ErrEditItem           error = errors.New("Failed to edit the item")
	ErrDeleteItem         error = errors.New("Failed to delete the item")

	ErrDeterminingUserIDFromJWT error = errors.New("Failed to determine UserID from jwt")
)
