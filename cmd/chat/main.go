package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	auth_jwt_service "github.com/egorkto/Chat-go/internal/auth/jwt/service"
	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	auth_jwt_transport_http "github.com/egorkto/Chat-go/internal/auth/jwt/transport/http"
	"github.com/egorkto/Chat-go/internal/logger"
	storage_postgres "github.com/egorkto/Chat-go/internal/storage/postgres"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	transport_http_echo_utils "github.com/egorkto/Chat-go/internal/transport/http/echo/utils"
	transport_http_server "github.com/egorkto/Chat-go/internal/transport/http/server"
	users_service "github.com/egorkto/Chat-go/internal/users/service"
	users_storage "github.com/egorkto/Chat-go/internal/users/storage"
	users_transport "github.com/egorkto/Chat-go/internal/users/transport/http"
	"github.com/egorkto/Chat-go/internal/validator"
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"

	_ "github.com/egorkto/Chat-go/docs"
)

// @title           Chat API
// @version         0.1
// @description     Chat go
// @contact.name   Egor
// @contact.url    http://github.com/egorkto
// @host      localhost:5845
// @BasePath  /
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
	dbCfg := storage_postgres.NewConfigMust()
	db, err := storage_postgres_gorm.New(dbCfg)
	if err != nil {
		logger.Error("initialize db: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Debug("Initializing echo router")
	e := transport_http_echo.NewRouter(logger.Logger)

	logger.Debug("Initializing jwt generator")
	jwtCfg := auth_jwt_token_manager.NewConfigMust()
	tokenManager, err := auth_jwt_token_manager.New(jwtCfg)
	if err != nil {
		logger.Error("initialize jwt generator: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	validator := validator.New()

	var authorizedRoutes []echo.Route
	var unauthorizedRoutes []echo.Route

	usersStorage := users_storage.New(db)
	usersService := users_service.New(usersStorage)
	usersTransport := users_transport.New(usersService)
	authorizedRoutes = append(authorizedRoutes, usersTransport.Routes()...)

	authService := auth_jwt_service.New(tokenManager, usersStorage, validator)
	authTransport := auth_jwt_transport_http.New(authService)
	unauthorizedRoutes = append(unauthorizedRoutes, authTransport.Routes()...)

	transport_http_echo_utils.WithMiddlewares(authorizedRoutes, tokenManager.EchoMiddleware())

	transport_http_echo_utils.AddMany(e, authorizedRoutes)
	transport_http_echo_utils.AddMany(e, unauthorizedRoutes)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	logger.Debug("Initializing HTTP server")
	serverCfg := transport_http_server.NewConfigMust()
	server := transport_http_server.New(serverCfg, e, logger.Logger)

	if err := server.Run(ctx); err != nil {
		logger.Error("server stopped", slog.String("error", err.Error()))
	}
}
