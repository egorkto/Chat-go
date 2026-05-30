package validator_test

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/egorkto/Chat-go/internal/validator"
	"github.com/stretchr/testify/require"
)

func TestValidatePassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid",
			password: gofakeit.Password(true, true, true, true, true, 10),
			wantErr:  false,
		},
		{
			name:     "Long password",
			password: gofakeit.Password(true, true, true, true, true, 120),
			wantErr:  true,
		},
		{
			name:     "Weak password",
			password: gofakeit.Password(true, false, true, false, false, 5),
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  true,
		},
	}

	v := validator.New()

	user := domain.NewUninitializedUser(
		gofakeit.Name(),
		gofakeit.Username(),
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.ValidatePassword(tc.password, user)
			if err != nil {
				if !errors.Is(err, domain.ErrValidationFailed) {
					t.Errorf(
						"unexpected error type: %s. Expected %s",
						err.Error(),
						domain.ErrValidationFailed.Error(),
					)
				}
			}

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
