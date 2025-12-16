# Porcelain Layer Guide (Do Not Touch Generated Code)

This doc explains how to build an ergonomic Go SDK on top of the generated clients under `internal/generated/**`. Keep generation as-is; only add porcelain in your own packages.

## Principles
- Never edit generated files. Compose/wrap them.
- Keep generation knobs in `config/*.yaml`; regenerate via `go generate ./...`.
- Prefer small, domain-focused helpers over giant facades.
- Centralize request wiring (base URL, auth, headers, timeouts, UA).
- Return Go-native types and clear errors; hide OpenAPI plumbing where possible.
- Accept `context.Context` on outward-facing methods.
- Prefer `WithResponse` variants for better status/body inspection; unwrap in porcelain.

## Where to Put Porcelain
- Use the existing `client` package as the entrypoint.
- Create subpackages/files per domain that match generated packages. Current minimal spec tags: `info` → `internal/generated/info`, `internal` → `internal/generated/internal`. As new tags are added, mirror them (e.g., `aws`, `pods`, etc.).
- Keep shared config/types in `client` (e.g., Config, option setters). Avoid exposing generated package names in the public API.

## Construction Pattern
Use one shared HTTP client and a single request editor; bind each generated client with it.
```go
type Client struct {
	Info     *InfoClient
	Internal *InternalClient
	// add more domains as tags/specs grow
}

type Option func(*Config)

func New(opts ...Option) (*Client, error) {
    cfg := DefaultConfig()
    for _, opt := range opts { opt(&cfg) }

    base, err := url.Parse(cfg.BaseURL)
    if err != nil { return nil, fmt.Errorf("base url: %w", err) }

    hc := pickHTTPClient(cfg) // set timeout; optionally wrap with retry transport
    editor := makeRequestEditor(cfg) // user-agent, auth header, api-key, extra headers

    ic, err := geninfo.NewClient(base.String(),
        geninfo.WithHTTPClient(hc),
        geninfo.WithRequestEditorFn(editor))
    if err != nil { return nil, fmt.Errorf("info client: %w", err) }

    inc, err := geninternal.NewClient(base.String(),
        geninternal.WithHTTPClient(hc),
        geninternal.WithRequestEditorFn(editor))
    if err != nil { return nil, fmt.Errorf("internal client: %w", err) }

    return &Client{
        Info:     &InfoClient{gen: ic},
        Internal: &InternalClient{gen: inc},
    }, nil
}
```
Option helpers to implement in porcelain (examples):
```go
func WithBaseURL(u string) Option        { return func(c *Config) { c.BaseURL = u } }
func WithHTTPClient(h *http.Client) Option { return func(c *Config) { c.HTTPClient = h } }
func WithUserAgent(ua string) Option     { return func(c *Config) { c.UserAgent = ua } }
func WithAuthHeader(v string) Option     { return func(c *Config) { c.AuthHeader = v } }
func WithAPIKey(v string) Option         { return func(c *Config) { c.APIKey = v } }
func WithRequestEditor(fn RequestEditor) Option {
    return func(c *Config) { c.RequestEditors = append(c.RequestEditors, fn) }
}
```
`makeRequestEditor` can chain `cfg.RequestEditors` plus built-ins.

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
    if err := expect(resp.HTTPResponse, http.StatusOK); err != nil {
        return nil, errWithBody(err, resp.Body)
    }
    if resp.JSON200 == nil {
        return nil, fmt.Errorf("list pods: empty body")
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
- Current minimal spec: `internal/generated/info` and `internal/generated/internal`, each already contains both client and types. Reuse those types; don’t duplicate models in porcelain.
- As new tags are added, new generated packages will include their own types—consume them directly or alias them in porcelain for a stable public API.
- Keep generated package names out of your public API for future flexibility; re-export or wrap types if you want a stable surface.

## Error Handling Pattern
Provide helpers to map `WithResponse` outputs to errors and keep bodies for debugging:
```go
func expect(r *http.Response, want int) error {
    if r.StatusCode == want {
        return nil
    }
    return fmt.Errorf("unexpected status %d", r.StatusCode)
}

func errWithBody(err error, body []byte) error {
    if len(body) == 0 {
        return err
    }
    return fmt.Errorf("%w: %s", err, string(body))
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

