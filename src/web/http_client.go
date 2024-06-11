package web

import (
	"net/http"
	"time"

	"github.com/andrepostiga/team-cron-notifier/src/config"
)

func NewClient(options config.HttpClientConfig) *http.Client {
	return &http.Client{
		Timeout: time.Duration(options.TimeoutInSeconds) * time.Second,
	}
}
