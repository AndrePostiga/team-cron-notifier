package web

import (
	"log/slog"
	"net/http"
	"time"
)

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
	Logger  *slog.Logger
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Log the HTTP request details
	lrt.Logger.Info("HTTP Request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.Scheme),
		slog.String("host", req.Host),
	)

	// Perform the actual request
	resp, err := lrt.Proxied.RoundTrip(req)
	duration := time.Since(start)

	// Determine if the request was successful
	isSuccessful := err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300

	// Log the HTTP response details and the IsSuccessful flag
	if err != nil {
		lrt.Logger.Error("HTTP Request failed",
			slog.String("method", req.Method),
			slog.String("url", req.URL.Host),
			slog.String("host", req.Host),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()),
			slog.Bool("IsSuccessful", isSuccessful),
		)
		return nil, err
	}

	lrt.Logger.Info("HTTP Response",
		slog.String("method", req.Method),
		slog.String("url", req.URL.Host),
		slog.String("host", req.Host),
		slog.Int("status", resp.StatusCode),
		slog.Duration("duration", duration),
		slog.Bool("IsSuccessful", isSuccessful),
	)

	return resp, nil
}
