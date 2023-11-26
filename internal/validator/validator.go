package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// Validator struct contains validation error messages for our form fields.
type Validator struct {
	FieldErrors map[string]string
}

// Valid checks if there are any validation errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldErorr adds a new validation error on to the validator struct
func (v *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't already
	// initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds a validation error only if the validation check is not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
