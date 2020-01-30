package validation

import (
	"reflect"
	"strings"
	"unicode"

	"gopkg.in/go-playground/validator.v9"
)

// ComplexityValidator validates that a string meets the minimum password
// complexity requirements.
var ComplexityValidator validator.Func = func(fl validator.FieldLevel) bool {
	pass := strings.TrimSpace(fl.Field().String())
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	hasMinLength := len(pass) > 10

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial && hasMinLength
}

// NotBlank is the validation function for validating if the current field
// has a value or length greater than zero, or is not a space only string.
var NotBlank validator.Func = func(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
