package auth_jwt_service

import (
	"fmt"
	"time"
)

func (s AuthService) GetTokenExpires(token string) (time.Time, error) {
	expires, err := s.tokenManager.GetExpires(token)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"getting expires: %w",
			err,
		)
	}
	return expires, nil
}
