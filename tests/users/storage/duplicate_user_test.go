package tests_users

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/egorkto/Chat-go/internal/domain"
	users_storage "github.com/egorkto/Chat-go/internal/users/storage"
	tests_postgres "github.com/egorkto/Chat-go/tests/storage/postgres"
	"github.com/stretchr/testify/require"
)

func TestDuplicateUser(t *testing.T) {
	db := tests_postgres.NewDB(10*time.Second, "migrations/init.sql", t)
	storage := users_storage.New(db)
	ctx := context.Background()

	t.Run("Test unique login constraint", func(t *testing.T) {
		user := domain.NewUser(0, 0, gofakeit.Name(), gofakeit.Username())
		password := gofakeit.Password(true, true, true, true, true, 10)

		_, err := storage.CreateUser(ctx, user, password)
		if err != nil {
			t.Fatalf("failed to create first user: %s", err.Error())
		}

		_, err = storage.CreateUser(ctx, user, password)
		require.Error(t, err)
	})
}
