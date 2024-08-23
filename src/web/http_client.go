package web

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/andrepostiga/team-cron-notifier/src/config"
)

func NewClient(logger *slog.Logger, options config.HttpClientConfig) *http.Client {
	return &http.Client{
		Timeout: time.Duration(options.TimeoutInSeconds) * time.Second,
		Transport: &LoggingRoundTripper{
			Proxied: http.DefaultTransport,
			Logger:  logger,
		},
	}
}
