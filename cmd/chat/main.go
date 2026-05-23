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
	chat_service "github.com/egorkto/Chat-go/internal/chat/service"
	chat_storage_postgres "github.com/egorkto/Chat-go/internal/chat/storage/postgres"
	chat_transport_http "github.com/egorkto/Chat-go/internal/chat/transport/http"
	chat_transport_websocket_hub "github.com/egorkto/Chat-go/internal/chat/transport/websocket/hub"
	"github.com/egorkto/Chat-go/internal/logger"
	storage_postgres "github.com/egorkto/Chat-go/internal/storage/postgres"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	transport_http_echo_pages "github.com/egorkto/Chat-go/internal/transport/http/echo/pages"
	transport_http_echo_utils "github.com/egorkto/Chat-go/internal/transport/http/echo/utils"
	transport_http_server "github.com/egorkto/Chat-go/internal/transport/http/server"
	users_service "github.com/egorkto/Chat-go/internal/users/service"
	users_storage_postgres "github.com/egorkto/Chat-go/internal/users/storage/postgres"
	users_transport "github.com/egorkto/Chat-go/internal/users/transport/http"
	"github.com/egorkto/Chat-go/internal/validator"
	"github.com/labstack/echo/v5"

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

	logger.Debug("Initializing echo router")
	e := transport_http_echo.NewRouter(logger.Logger)

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
	wsHub := chat_transport_websocket_hub.New(chatService, e.Logger, wsHubCfg)
	chatTransport := chat_transport_http.New(chatService, wsHub)

	logger.Debug("Registration routes")
	var authorizedRoutes []echo.Route
	var unauthorizedRoutes []echo.Route

	unauthorizedRoutes = append(unauthorizedRoutes, authTransport.Routes()...)
	unauthorizedRoutes = append(unauthorizedRoutes, transport_http_echo_pages.Routes()...)
	authorizedRoutes = append(authorizedRoutes, usersTransport.Routes()...)
	authorizedRoutes = append(authorizedRoutes, chatTransport.Routes()...)

	transport_http_echo_utils.WithMiddlewares(authorizedRoutes, tokenManager.HeaderMiddleware())

	wsRoute := chatTransport.ConnectWebsocketRoute()
	wsRoute.Middlewares = append(wsRoute.Middlewares, tokenManager.QueryMiddleware())
	authorizedRoutes = append(authorizedRoutes, wsRoute)

	transport_http_echo_utils.AddMany(e, authorizedRoutes)
	transport_http_echo_utils.AddMany(e, unauthorizedRoutes)

	logger.Debug("Initializing HTTP server")
	serverCfg := transport_http_server.NewConfigMust()
	server := transport_http_server.New(serverCfg, e, logger.Logger)

	go wsHub.StartBroadcasting()

	if err := server.Run(ctx); err != nil {
		logger.Error("server stopped", slog.String("error", err.Error()))
	}
}
