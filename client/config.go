package client

import (
	"net/http"
	"time"
)

const (
	defaultBaseURL = "http://localhost.localstack.cloud:4566"
	defaultTimeout = 60 * time.Second
)

// Config holds client configuration for the LocalStack SDK.
type Config struct {
	// BaseURL points to the LocalStack endpoint.
	BaseURL string
	// HTTPClient allows callers to supply a custom HTTP client (e.g., with tracing or retries).
	HTTPClient *http.Client
}

// DefaultConfig returns a Config pre-populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL: defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}
