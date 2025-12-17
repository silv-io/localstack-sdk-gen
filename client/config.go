package client

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "http://localhost.localstack.cloud:4566"
	defaultTimeout = 60 * time.Second
	defaultUserAgent = "localstack-sdk-go"
)

// Config holds client configuration for the LocalStack SDK.
type Config struct {
	// BaseURL points to the LocalStack endpoint.
	BaseURL string
	// HTTPClient allows callers to supply a custom HTTP client (e.g., with tracing or retries).
	HTTPClient *http.Client
	// Timeout configures the timeout used when HTTPClient is nil.
	Timeout time.Duration
	// UserAgent overrides the User-Agent header if non-empty.
	UserAgent string
	// AuthHeader sets the Authorization header value if non-empty.
	AuthHeader string
	// APIKey sets the x-api-key header value if non-empty.
	APIKey string
	// Headers are applied to every request (after built-ins).
	Headers http.Header
	// RequestEditors are applied after built-ins and Headers.
	RequestEditors []RequestEditor
}

// DefaultConfig returns a Config pre-populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL: defaultBaseURL,
		Timeout:  defaultTimeout,
		UserAgent: defaultUserAgent,
		Headers:  make(http.Header),
	}
}

// Option mutates Config during construction.
type Option func(*Config)

// RequestEditor mutates an outgoing request before it is sent.
type RequestEditor func(ctx context.Context, req *http.Request) error

func WithBaseURL(u string) Option { return func(c *Config) { c.BaseURL = u } }

func WithHTTPClient(h *http.Client) Option { return func(c *Config) { c.HTTPClient = h } }

func WithTimeout(d time.Duration) Option { return func(c *Config) { c.Timeout = d } }

func WithUserAgent(ua string) Option { return func(c *Config) { c.UserAgent = ua } }

func WithAuthHeader(v string) Option { return func(c *Config) { c.AuthHeader = v } }

func WithAPIKey(v string) Option { return func(c *Config) { c.APIKey = v } }

// WithHeader sets a single header value, overwriting any existing values.
func WithHeader(key, value string) Option {
	return func(c *Config) {
		if c.Headers == nil {
			c.Headers = make(http.Header)
		}
		c.Headers.Set(key, value)
	}
}

// WithHeaders merges all headers into the Config, overwriting existing keys.
func WithHeaders(h http.Header) Option {
	return func(c *Config) {
		if len(h) == 0 {
			return
		}
		if c.Headers == nil {
			c.Headers = make(http.Header)
		}
		for k, vs := range h {
			c.Headers.Del(k)
			for _, v := range vs {
				c.Headers.Add(k, v)
			}
		}
	}
}

func WithRequestEditor(fn RequestEditor) Option {
	return func(c *Config) { c.RequestEditors = append(c.RequestEditors, fn) }
}
