package domain_test

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestDomain_UserValidation(t *testing.T) {
	testCases := []struct {
		name    string
		user    domain.User
		wantErr bool
	}{
		{
			name: "Valid",
			user: domain.NewUninitializedUser(
				gofakeit.Name(),
				gofakeit.Username(),
			),
			wantErr: false,
		},
		{
			name: "Short full name",
			user: domain.NewUninitializedUser(
				gofakeit.LetterN(2),
				gofakeit.Username(),
			),
			wantErr: true,
		},
		{
			name: "Long full name",
			user: domain.NewUninitializedUser(
				gofakeit.LetterN(110),
				gofakeit.Username(),
			),
			wantErr: true,
		},
		{
			name: "Empty full name",
			user: domain.NewUninitializedUser(
				"",
				gofakeit.Username(),
			),
			wantErr: true,
		},
		{
			name: "Short login",
			user: domain.NewUninitializedUser(
				gofakeit.Name(),
				gofakeit.LetterN(2),
			),
			wantErr: true,
		},
		{
			name: "Long login",
			user: domain.NewUninitializedUser(
				gofakeit.Name(),
				gofakeit.LetterN(30),
			),
			wantErr: true,
		},
		{
			name: "Empty login",
			user: domain.NewUninitializedUser(
				gofakeit.Name(),
				"",
			),
			wantErr: true,
		},
		{
			name: "Empty",
			user: domain.NewUninitializedUser(
				"",
				"",
			),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.user.Validate()
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
