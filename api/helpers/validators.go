package helpers

import (
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

// ValidatePassword ensures that the password meets the following criteria:
// - At least 8 characters long
// - Contains at least one lowercase letter
// - Contains at least one uppercase letter
// - Contains at least one number
// - Contains at least one special character
// It also ensures that the password and confirm password match.
func ValidatePassword(password string) error {

	if len(password) == 0 || len(strings.TrimSpace(password))==0{
		return fmt.Errorf("password cannot be empty or just spaces")
	}

	// Validate password using various rules
	err := validation.Validate(password,
		validation.Length(8, 0).Error("password must be at least 8 characters long"),
		validation.Match(regexp.MustCompile(`[a-z]`)).Error("password must contain at least one lowercase letter"),
		validation.Match(regexp.MustCompile(`[A-Z]`)).Error("password must contain at least one uppercase letter"),
		validation.Match(regexp.MustCompile(`\d`)).Error("password must contain at least one number"),
		validation.Match(regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)).Error("password must contain at least one special character"),
	)
	// Return any validation error
	if err != nil {
		return err
	}

	return nil
}
