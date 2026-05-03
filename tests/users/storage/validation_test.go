package tests_users

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/egorkto/Chat-go/internal/domain"
	users_storage "github.com/egorkto/Chat-go/internal/users/storage"
	tests_postgres "github.com/egorkto/Chat-go/tests/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	testCases := []struct {
		name     string
		user     domain.User
		password string
		wantErr  bool
	}{
		{
			name:     "Valid data",
			user:     domain.NewUser(0, 0, gofakeit.Name(), gofakeit.Username()),
			password: gofakeit.Password(true, true, true, true, true, 10),
			wantErr:  false,
		},
		{
			name:     "Empty password",
			user:     domain.NewUser(0, 0, gofakeit.Name(), gofakeit.Username()),
			password: "",
			wantErr:  true,
		},
		{
			name:     "Short full name",
			user:     domain.NewUser(0, 0, gofakeit.LetterN(2), gofakeit.Username()),
			password: gofakeit.Password(true, true, true, true, true, 10),
			wantErr:  true,
		},
		{
			name:     "Long full name",
			user:     domain.NewUser(0, 0, gofakeit.LetterN(21), gofakeit.Username()),
			password: gofakeit.Password(true, true, true, true, true, 10),
			wantErr:  true,
		},
	}

	db := tests_postgres.NewDB(10*time.Second, "migrations/init.sql", t)
	storage := users_storage.New(db)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := storage.CreateUser(ctx, tc.user, tc.password)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
