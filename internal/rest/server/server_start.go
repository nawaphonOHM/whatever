package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/nawaphonOHM/whatever/v2/pkg/logging"
)

// listenAndServeAsync starts ListenAndServe and reports non-close errors.
func listenAndServeAsync(s *Server) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()
	return errChan
}

// gracefulShutdown shuts down the HTTP server using ShutdownTimeout.
func gracefulShutdown(ctx context.Context, s *Server) error {
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		s.Config.ShutdownTimeout,
	)
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		logging.ErrorContext(ctx, "REST server graceful shutdown failed", "error", err.Error())
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}
	logging.InfoContext(ctx, "REST server shut down successfully")
	return nil
}

// handleServeErr logs serve failure and wraps the non-nil error.
func handleServeErr(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	logging.ErrorContext(ctx, "REST server failed to start", "error", err.Error())
	return fmt.Errorf("http server failed to start: %w", err)
}

// waitForServerExit blocks until serve fails or notify context ends.
func waitForServerExit(
	ctx context.Context,
	errChan <-chan error,
	s *Server,
) error {
	select {
	case err := <-errChan:
		return handleServeErr(ctx, err)
	case <-ctx.Done():
		logging.InfoContext(
			ctx,
			"shutting down REST server gracefully",
			"shutdown_timeout", s.Config.ShutdownTimeout.String(),
		)
		return nil
	}
}

// Start launches the HTTP server until cancel or OS signal.
func (s *Server) Start(ctx context.Context) error {
	notifyCtx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	logging.InfoContext(
		notifyCtx,
		"starting REST server",
		"service", s.Config.DisplayName,
		"version", s.Config.AppVersion,
		"addr", s.httpServer.Addr,
		"host", s.Config.Host,
		"port", s.Config.Port,
		"mode", s.Config.Mode,
		"shutdown_timeout", s.Config.ShutdownTimeout.String(),
		"access_log", s.Config.EnableAccessLog,
		"metrics", s.Config.EnableMetrics,
	)

	errChan := listenAndServeAsync(s)
	if err := waitForServerExit(notifyCtx, errChan, s); err != nil {
		return err
	}
	return gracefulShutdown(ctx, s)
}

// Shutdown initiates graceful shutdown using the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
