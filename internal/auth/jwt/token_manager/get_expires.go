package auth_jwt_token_manager

import (
	"fmt"
	"time"
)

func (tm TokenManager) GetExpires(token string) (time.Time, error) {
	jwt, err := tm.Verify(token)
	if err != nil {
		return time.Time{}, fmt.Errorf("verifying token: %w", err)
	}

	date, err := jwt.Claims.GetExpirationTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("getting expiration time: %w", err)
	}

	return date.Time, nil
}
