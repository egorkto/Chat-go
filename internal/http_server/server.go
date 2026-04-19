package http_server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type HTTPServer struct {
	config  Config
	handler http.Handler
	log     *slog.Logger
}

func New(
	config Config,
	handler http.Handler,
	log *slog.Logger,
) *HTTPServer {
	return &HTTPServer{
		config:  config,
		handler: handler,
		log:     log,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%v", s.config.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: s.handler,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", slog.Int("port", s.config.Port))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
