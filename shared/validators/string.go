package validators

import (
	"fmt"
	"strings"
)

type StringValidatorFunc func(string, string) error

func NonEmptyString() StringValidatorFunc {
	return func(value string, name string) error {
		if value == "" {
			return fmt.Errorf("'%s' value is empty", name)
		}
		return nil
	}
}

func StringWithoutSpaces() StringValidatorFunc {
	return func(value string, name string) error {
		if strings.ContainsAny(value, " \t\n") {
			return fmt.Errorf("'%s' value should not contain spaces", name)
		}
		return nil
	}
}
