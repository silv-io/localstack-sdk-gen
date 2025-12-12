# Porcelain Layer Guide (Do Not Touch Generated Code)

This doc explains how to build an ergonomic Go SDK on top of the generated clients under `internal/generated/**`. Keep generation as-is; only add porcelain in your own packages.

## Principles
- Never edit generated files. Compose/wrap them.
- Prefer small, domain-focused helpers over giant facades.
- Keep request wiring (base URL, auth, headers, timeouts) centralized.
- Return Go-native types and clear errors; hide OpenAPI plumbing where possible.
- Accept `context.Context` on outward-facing methods.

## Where to Put Porcelain
- Use the existing `client` package as the entrypoint.
- Create subpackages or files per domain (e.g., `client/localstack`, `client/aws`, `client/pods`, `client/chaos`, `client/replicator`) that wrap the generated packages `genlocalstack`, `genaws`, `genpods`, `genchaos`, `genrepl`.
- Keep shared config/types in `client` (e.g., Config, option setters).

## Construction Pattern
```go
type Client struct {
	Localstack *LocalstackClient
	AWS        *AWSClient
	Pods       *PodsClient
	Chaos      *ChaosClient
	Replicator *ReplicatorClient
}

type Option func(*Config)

func New(opts ...Option) (*Client, error) {
    cfg := DefaultConfig()
    for _, opt := range opts { opt(&cfg) }

    base, err := url.Parse(cfg.BaseURL)
    if err != nil { return nil, fmt.Errorf("base url: %w", err) }

    // one shared http.Client (with timeouts, tracing, retries if you add them)
    hc := cfg.HTTPClient
    if hc == nil { hc = defaultHTTPClient() }

    // request editor for auth/headers
    editor := func(ctx context.Context, req *http.Request) error {
        req.Header.Set("User-Agent", cfg.UserAgent)
        if cfg.AuthHeader != "" {
            req.Header.Set("Authorization", cfg.AuthHeader)
        }
        if cfg.APIKey != "" {
            req.Header.Set("x-api-key", cfg.APIKey)
        }
        return nil
    }

    lc, err := genlocalstack.NewClient(base.String(), genlocalstack.WithHTTPClient(hc), genlocalstack.WithRequestEditorFn(editor))
    if err != nil { return nil, fmt.Errorf("localstack client: %w", err) }
    // repeat for other generated clients...

    return &Client{
        Localstack: &LocalstackClient{gen: lc},
        AWS:        &AWSClient{gen: ac},
        Pods:       &PodsClient{gen: pc},
        Chaos:      &ChaosClient{gen: cc},
        Replicator: &ReplicatorClient{gen: rc},
    }, nil
}
```

## Porcelain Methods
Wrap generated calls with:
- Context-first signatures: `func (c *PodsClient) List(ctx context.Context) ([]Pod, error)`
- Friendly inputs: use Go types instead of pointers to scalars; set defaults internally.
- Friendly outputs: unwrap `*_WithResponse` results to `(T, error)`.
- Error wrapping: normalize HTTP errors into Go errors with status/operation info.
- Pagination helpers: provide iterators or “ListAll” helpers that handle tokens/limits.
- Idempotent convenience: e.g., “Upsert” built from generated Get+Create/Update.

Example:
```go
func (c *PodsClient) List(ctx context.Context) ([]genpods.PodSummary, error) {
    resp, err := c.gen.ListPodsWithResponse(ctx)
    if err != nil {
        return nil, fmt.Errorf("list pods: %w", err)
    }
    if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
        return nil, fmt.Errorf("list pods: unexpected status %d", resp.StatusCode())
    }
    return *resp.JSON200, nil
}
```

## Options to Support
- Base URL override.
- Custom `http.Client`.
- Request editor hooks (auth headers, tracing, extra headers).
- User-Agent override.
- Timeouts/retries (via the provided http.Client; you can add a retry transport if desired).

## Generated Package Use
- Types per domain live in:
  - `internal/generated/types` (shared models)
  - `internal/generated/localstack`
  - `internal/generated/aws`
  - `internal/generated/pods`
  - `internal/generated/chaos`
  - `internal/generated/replicator`
- Prefer reusing shared types when sensible; otherwise wrap/alias in porcelain if more ergonomic.

## Error Handling Pattern
Provide a small helper to map `WithResponse` outputs to errors:
```go
func check(resp interface{ StatusCode() int }, okStatus int, body []byte) error {
    if resp.StatusCode() == okStatus {
        return nil
    }
    return fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(body))
}
```
Then use it in helpers after calling `*_WithResponse`.

## Testing Guidance (for your future additions)
- Unit-test porcelain methods with a fake `http.RoundTripper`.
- Add golden tests for error wrapping and pagination helpers.
- Keep generated code out of tests; test through the porcelain API.

## Generation Hygiene
- To regenerate: `go generate ./...` (uses configs in `config/`).
- Do not hand-edit anything under `internal/generated/**`.
- If you add configs or options, keep them in `config/` and reference via `go:generate`.

