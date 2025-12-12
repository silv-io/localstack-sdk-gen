package client

import (
	"net/http"
)

// Client is the porcelain entrypoint. It will wrap generated clients and expose
// ergonomic helpers in the future.
type Client struct {
	cfg Config
}

// New constructs a Client with defaults applied. This is intentionally light-weight
// and does not yet wire generated code; that can be added as porcelain evolves.
func New(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{cfg: cfg}
}
