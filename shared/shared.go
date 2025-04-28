package shared

import (
	"errors"
	"fmt"
	"os"
	"vmm/shared/validators"
)

func MustValidateName(name string) {
	Must(ValidateName(name))
}
func ValidateName(name string) error {
	return ValidateString(name, "name", validators.NonEmptyString(), validators.StringWithoutSpaces())
}

func MustValidateString(value string, name string) {
	Must(ValidateString(value, name, validators.NonEmptyString()))
}

func ValidateString(value string, name string, validators ...validators.StringValidatorFunc) error {
	var validationErrors error
	for _, v := range validators {
		validationErrors = errors.Join(validationErrors, v(value, name))
	}
	return validationErrors
}

func Must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}
