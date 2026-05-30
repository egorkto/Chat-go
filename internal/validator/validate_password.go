package validator

import (
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/nbutton23/zxcvbn-go"
)

func (v Validator) ValidatePassword(password string, user domain.User) error {
	lenPass := len([]rune(password))
	if lenPass > 100 || lenPass == 0 {
		return domain.NewValidationError(map[string]string{
			"password": "length must be between 0 and 101",
		})
	}

	match := zxcvbn.PasswordStrength(password, []string{user.FullName(), user.Login()})
	if match.Score < 2 {
		return domain.NewValidationError(map[string]string{
			"password": "password is too weak",
		})
	}

	return nil
}
