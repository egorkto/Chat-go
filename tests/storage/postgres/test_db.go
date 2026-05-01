package tests_postgres

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	db_gorm_postgres "github.com/egorkto/Chat-go/internal/db/gorm/postgres"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestDB struct {
	*gorm.DB
	timeout time.Duration
}

func NewDB(timeout time.Duration, t *testing.T) TestDB {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %s", err.Error())
	}

	initSql := wd + "/init.sql"

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:17.9-alpine",
		postgres.WithInitScripts(initSql),
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
	return db_gorm_postgres.MapError(err)
}
