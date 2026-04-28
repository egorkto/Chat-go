package validator

import (
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/nbutton23/zxcvbn-go"
)

func (v Validator) ValidatePassword(password string, fullName string) error {
	lenPass := len([]byte(password))
	if lenPass > 100 || lenPass == 0 {
		return fmt.Errorf(
			"password length must be between 0 and 101: %w",
			domain.ErrInvalidArgument,
		)
	}

	match := zxcvbn.PasswordStrength(password, []string{fullName})
	if match.Score < 2 {
		return fmt.Errorf(
			"password is too weak: %w",
			domain.ErrInvalidArgument,
		)
	}
	return nil
}
