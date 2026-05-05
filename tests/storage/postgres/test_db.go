package tests_postgres

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
	tests_utils "github.com/egorkto/Chat-go/tests/utils"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestDB struct {
	*gorm.DB
	timeout time.Duration
}

func NewDB(timeout time.Duration, initMigration string, t *testing.T) TestDB {
	ctx := context.Background()

	var migrations []string

	migrationsDir := filepath.Join(tests_utils.GetProjectRoot(), "migrations")
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

	db, err := gorm.Open(gorm_postgres.Open(connStr))
	if err != nil {
		t.Fatalf("failed to open gorm db: %s", err.Error())
	}

	return TestDB{
		DB:      db,
		timeout: timeout,
	}
}

func (t TestDB) WithTimeoutContext(ctx context.Context) context.CancelFunc {
	timeoutCtx, cancel := context.WithTimeout(ctx, t.timeout)

	t.DB = t.DB.WithContext(timeoutCtx)

	return cancel
}

func (t TestDB) Create(dest interface{}) error {
	result := t.DB.Create(dest)
	err := result.Error
	if err != nil {
		log.Printf("[TestDB] Create error: %v", err)
	}
	return storage_postgres_gorm.MapError(err)
}

func (t TestDB) First(dest interface{}, query interface{}, args ...interface{}) error {
	result := t.DB.Where(query, args...).First(dest)
	err := result.Error
	if err != nil {
		log.Printf("[TestDB] First error: %v", err)
	}
	return storage_postgres_gorm.MapError(err)
}
