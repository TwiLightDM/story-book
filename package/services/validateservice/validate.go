package validateservice

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile(
	`^[a-zA-Z0-9._%+-]+@([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`,
)

type ValidationService interface {
	IsValidEmail(email string) error
	IsStrongPassword(password string) error
}

type validationService struct {
	minPasswordSize int
}

func NewValidationService(minPasswordSize int) ValidationService {
	return &validationService{minPasswordSize: minPasswordSize}
}

func (v *validationService) IsValidEmail(email string) error {
	if emailRegex.MatchString(email) {
		return nil
	}

	return ErrBadEmail
}

func (v *validationService) IsStrongPassword(password string) error {
	if len(password) < v.minPasswordSize {
		return ErrBadPassword
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if hasUpper && hasLower && hasDigit {
		return nil
	}

	return ErrBadPassword
}
