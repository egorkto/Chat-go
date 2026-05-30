package tests_postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
	"github.com/egorkto/Chat-go/internal/utils"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(
	timeout time.Duration,
	t *testing.T,
) (*storage_postgres_gorm.GormDB, error) {
	ctx := context.Background()

	var migrations []string

	migrationsDir := filepath.Join(utils.GetProjectRoot(), "migrations")
	entries, err := os.ReadDir(migrationsDir)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			migrations = append(migrations, filepath.Join(migrationsDir, entry.Name()))
		}
	}

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:17.9-alpine",
		postgres.WithInitScripts(migrations...),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err.Error())
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err.Error())
	}

	db, err := gorm.Open(gorm_postgres.Open(connStr), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("failed to open gorm db: %s", err.Error())
	}

	gormDB := storage_postgres_gorm.GormDB{
		DB: db,
	}

	if err := gormDB.SetTimeout(timeout); err != nil {
		return nil, fmt.Errorf("set timeout: %w", err)
	}

	return &gormDB, nil
}
