package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func BuildApi(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	// Start server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()
	log.Println("Server started on :3000")

	// Wait for interrupt signal to gracefully shutdown the server
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("Context canceled, shutting down server.")
	case sig := <-shutdown:
		log.Printf("Received signal: %s, shutting down server.", sig)
	}

	// Create a context with timeout for the shutdown
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
