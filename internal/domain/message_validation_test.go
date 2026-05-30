package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestDomain_MessageValidation(t *testing.T) {
	validUser := domain.NewUninitializedUser(
		gofakeit.Name(),
		gofakeit.Username(),
	)

	invalidUser := domain.NewUninitializedUser(
		"",
		"",
	)

	testCases := []struct {
		name    string
		message domain.Message
		wantErr bool
	}{
		{
			name: "Valid",
			message: domain.NewUninitializedMessage(
				validUser,
				gofakeit.LetterN(10),
				time.Now(),
			),
			wantErr: false,
		},
		{
			name: "Invalid sender",
			message: domain.NewUninitializedMessage(
				invalidUser,
				gofakeit.LetterN(10),
				time.Now(),
			),
			wantErr: true,
		},
		{
			name: "Empty text",
			message: domain.NewUninitializedMessage(
				validUser,
				"",
				time.Now(),
			),
			wantErr: true,
		},
		{
			name: "Long text",
			message: domain.NewUninitializedMessage(
				validUser,
				gofakeit.LetterN(5000),
				time.Now(),
			),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.message.Validate()
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
