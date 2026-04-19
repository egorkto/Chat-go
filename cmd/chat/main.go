package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/egorkto/Chat-go/internal/auth"
	auth_service "github.com/egorkto/Chat-go/internal/auth/service"
	auth_transport "github.com/egorkto/Chat-go/internal/auth/transport/http"
	"github.com/egorkto/Chat-go/internal/db"
	db_gorm_postgres "github.com/egorkto/Chat-go/internal/db/gorm/postgres"
	echo_utils "github.com/egorkto/Chat-go/internal/echo/utils.go"
	"github.com/egorkto/Chat-go/internal/http_server"
	"github.com/egorkto/Chat-go/internal/logger"
	users_storage "github.com/egorkto/Chat-go/internal/users/storage"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/presbrey/pkg/echovalidator"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	loggerCfg := logger.NewConfigMust()
	logger, err := logger.New(loggerCfg)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to initialize logger: %w", err))
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Initializing postgres db")
	dbCfg := db.NewConfigMust()
	db, err := db_gorm_postgres.New(dbCfg)
	if err != nil {
		logger.Error("initialize db: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Debug("Initializing echo router")
	e := echo.New()
	e.Logger = logger.Logger
	e.Validator = echovalidator.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	logger.Debug("Initializing jwt generator")
	jwtCfg := auth.NewJWTConfigMust()
	jwtGenerator, err := auth.NewJWTGenerator(jwtCfg)
	if err != nil {
		logger.Error("initialize jwt generator: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	usersStorage := users_storage.New(db)
	authService := auth_service.New(jwtGenerator, usersStorage)
	authTransport := auth_transport.New(authService)

	echo_utils.AddMany(e, authTransport.Routes())

	logger.Debug("Initializing HTTP server")
	serverCfg := http_server.NewConfigMust()
	server := http_server.New(serverCfg, e, logger.Logger)

	if err := server.Run(ctx); err != nil {
		logger.Error("server stopped", slog.String("error", err.Error()))
	}
}
