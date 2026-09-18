package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
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
func gracefulShutdown(s *Server) error {
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		s.Config.ShutdownTimeout,
	)
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}
	return nil
}

// mapServeErr wraps a non-nil serve error.
func mapServeErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("http server failed to start: %w", err)
}

// waitForServerExit blocks until serve fails or notify context ends.
func waitForServerExit(
	ctx context.Context,
	errChan <-chan error,
) error {
	select {
	case err := <-errChan:
		return mapServeErr(err)
	case <-ctx.Done():
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

	errChan := listenAndServeAsync(s)
	if err := waitForServerExit(notifyCtx, errChan); err != nil {
		return err
	}
	return gracefulShutdown(s)
}

// Shutdown initiates graceful shutdown using the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
