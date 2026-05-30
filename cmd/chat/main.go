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
	auth_jwt_transport_authorizer "github.com/egorkto/Chat-go/internal/auth/jwt/transport/authorizer"
	auth_jwt_transport_http "github.com/egorkto/Chat-go/internal/auth/jwt/transport/http"
	chat_service "github.com/egorkto/Chat-go/internal/chat/service"
	chat_storage_postgres "github.com/egorkto/Chat-go/internal/chat/storage/postgres"
	chat_transport_http "github.com/egorkto/Chat-go/internal/chat/transport/http"
	chat_transport_websocket_hub "github.com/egorkto/Chat-go/internal/chat/transport/websocket/hub"
	"github.com/egorkto/Chat-go/internal/logger"
	pages_transport "github.com/egorkto/Chat-go/internal/pages/transport"
	storage_postgres "github.com/egorkto/Chat-go/internal/storage/postgres"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	transport_http_echo_router "github.com/egorkto/Chat-go/internal/transport/http/echo/router"
	transport_http_server "github.com/egorkto/Chat-go/internal/transport/http/server"
	users_service "github.com/egorkto/Chat-go/internal/users/service"
	users_storage_postgres "github.com/egorkto/Chat-go/internal/users/storage/postgres"
	users_transport "github.com/egorkto/Chat-go/internal/users/transport/http"
	"github.com/egorkto/Chat-go/internal/validator"

	_ "github.com/egorkto/Chat-go/docs"
)

// @title           Chat API
// @version         0.1
// @description     Chat go
// @contact.name   Egor
// @contact.url    http://github.com/egorkto
// @host      localhost:5845
// @BasePath  /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Вставьте токен в формате: Bearer <ваш_токен>
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

	logger.Debug("Initializing jwt token manager")
	jwtCfg := auth_jwt_token_manager.NewConfigMust()
	tokenManager, err := auth_jwt_token_manager.New(jwtCfg)
	if err != nil {
		logger.Error("new jwt token manager: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	validator := validator.New()

	logger.Debug("Initializing users")
	usersStorage := users_storage_postgres.New(db)
	usersService := users_service.New(usersStorage)
	usersTransport := users_transport.New(usersService)

	logger.Debug("Initializing auth")
	authService := auth_jwt_service.New(tokenManager, usersStorage, validator)
	authTransport := auth_jwt_transport_http.New(authService)

	logger.Debug("Initializing chat")
	msgStorage := chat_storage_postgres.New(db)
	chatService := chat_service.New(usersStorage, msgStorage)
	wsHubCfg := chat_transport_websocket_hub.NewConfigMust()
	wsHub := chat_transport_websocket_hub.New(wsHubCfg, chatService, logger.Logger)
	chatTransport := chat_transport_http.New(chatService, wsHub)

	logger.Debug("Initializing pages transport")
	pagesTransport := pages_transport.New()

	authorizer := auth_jwt_transport_authorizer.New(tokenManager)

	mc := transport_http.MiddlewaresContainer{
		HeaderAuth: authorizer.AuthorizeHeaderMiddleware(),
		QueryAuth:  authorizer.AuthorizeQueryMiddleware(),
	}

	logger.Debug("Initializing echo router")
	router := transport_http_echo_router.New(mc)
	echoRouter := router.NewRouter(logger.Logger, []transport_http_echo_router.HTTPHandler{
		usersTransport,
		authTransport,
		chatTransport,
		pagesTransport,
	})

	logger.Debug("Initializing HTTP server")
	serverCfg := transport_http_server.NewConfigMust()
	server := transport_http_server.New(serverCfg, echoRouter, logger.Logger)

	go wsHub.StartBroadcasting()

	if err := server.Run(ctx); err != nil {
		logger.Error("server stopped", slog.String("error", err.Error()))
	}
}
