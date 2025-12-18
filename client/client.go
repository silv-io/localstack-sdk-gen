package client

import (
	"fmt"
	"net/http"
	"net/url"

	awsapi "github.com/localstack/localstack-sdk-go/internal/generated/awsapi"
	chaosapi "github.com/localstack/localstack-sdk-go/internal/generated/chaosapi"
	infoapi "github.com/localstack/localstack-sdk-go/internal/generated/infoapi"
	internalapi "github.com/localstack/localstack-sdk-go/internal/generated/internalapi"
)

// Client is the porcelain entrypoint. It wraps generated clients and exposes
// ergonomic helpers.
type Client struct {
	cfg      Config
	Info     *InfoClient
	Internal *InternalClient
	Chaos    *ChaosClient
	AWS      *AWSClient
}

// New constructs a Client with defaults applied and wires all domain clients.
func New(opts ...Option) (*Client, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return NewWithConfig(cfg)
}

// NewWithConfig constructs a Client from an explicit Config.
func NewWithConfig(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}

	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("base url: %w", err)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = defaultTimeout
		}
		httpClient = &http.Client{Timeout: timeout}
	}

	editor := makeRequestEditor(cfg)

	infoGen, err := infoapi.NewClientWithResponses(baseURL.String(),
		infoapi.WithHTTPClient(httpClient),
		infoapi.WithRequestEditorFn(infoapi.RequestEditorFn(editor)),
	)
	if err != nil {
		return nil, fmt.Errorf("info client: %w", err)
	}

	internalGen, err := internalapi.NewClientWithResponses(baseURL.String(),
		internalapi.WithHTTPClient(httpClient),
		internalapi.WithRequestEditorFn(internalapi.RequestEditorFn(editor)),
	)
	if err != nil {
		return nil, fmt.Errorf("internal client: %w", err)
	}

	chaosGen, err := chaosapi.NewClientWithResponses(baseURL.String(),
		chaosapi.WithHTTPClient(httpClient),
		chaosapi.WithRequestEditorFn(chaosapi.RequestEditorFn(editor)),
	)
	if err != nil {
		return nil, fmt.Errorf("chaos client: %w", err)
	}

	awsGen, err := awsapi.NewClientWithResponses(baseURL.String(),
		awsapi.WithHTTPClient(httpClient),
		awsapi.WithRequestEditorFn(awsapi.RequestEditorFn(editor)),
	)
	if err != nil {
		return nil, fmt.Errorf("aws client: %w", err)
	}

	return &Client{
		cfg:      cfg,
		Info:     &InfoClient{gen: infoGen},
		Internal: &InternalClient{gen: internalGen},
		Chaos:    &ChaosClient{gen: chaosGen},
		AWS:      &AWSClient{gen: awsGen},
	}, nil
}
