package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error {
	if len(value) < minLength || len(value) > maxLength {
		return fmt.Errorf("must be between [%d-%d] characters", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 4, 12); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("must containe only letters, digits or underscores")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 24)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("must be a valid email address")
	}

	return nil
}
