package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/andrepostiga/team-cron-notifier/src/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func BuildApi(ctx context.Context) {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to load config: %v", err))
	}

	logLevel := new(slog.LevelVar)
	err = logLevel.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		slog.Error("Failed to load config", "error", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	mux := http.NewServeMux()

	// Health check handler with request and response logging
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Extract only the headers of interest
		headers := map[string]string{
			"Host":            r.Header.Get("Host"),
			"Accept-Encoding": r.Header.Get("Accept-Encoding"),
			"Content-Length":  r.Header.Get("Content-Length"),
			"Content-Type":    r.Header.Get("Content-Type"),
		}

		// Process the request
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

		// Log the request and response together
		slog.Info("Request and Response",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
			"headers", headers,
			"status", http.StatusOK,
			"duration", time.Since(start),
		)
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	// Start server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", "error", err)
		}
	}()
	slog.Info("Server started", "address", ":3000")

	// Wait for interrupt signal to gracefully shutdown the server
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		slog.Info("Context canceled, shutting down server.")
	case sig := <-shutdown:
		slog.Info("Received signal, shutting down server", "signal", sig)
	}

	// Create a context with timeout for the shutdown
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
