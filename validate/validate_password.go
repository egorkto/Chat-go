package validate

import (
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/nbutton23/zxcvbn-go"
)

func ValidatePassword(password string, fullName string) error {
	match := zxcvbn.PasswordStrength(password, []string{fullName})
	if match.Score < 2 {
		return fmt.Errorf(
			"validating password: %s: %w",
			"password is too weak",
			domain.ErrInvalidArgument,
		)
	}
	return nil
}
