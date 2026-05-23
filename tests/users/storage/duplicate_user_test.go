package tests_users

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/egorkto/Chat-go/internal/domain"
	users_storage_postgres "github.com/egorkto/Chat-go/internal/users/storage/postgres"
	tests_postgres "github.com/egorkto/Chat-go/tests/storage/postgres"
	"github.com/stretchr/testify/require"
)

func TestDuplicateUser(t *testing.T) {
	db, err := tests_postgres.NewDB(10*time.Second, "migrations/init.sql", t)
	if err != nil {
		t.Fatalf("new test db: %w", err)
	}
	storage := users_storage_postgres.New(db)
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
