# LocalStack SDK for Go

This project sets up Go SDK generation from the LocalStack OpenAPI spec using `oapi-codegen` via `go generate`. Generated code lives in `internal/generated` and is intended to be wrapped by the porcelain entrypoint in `client.go`.

## Layout
- `go.mod`: module `github.com/localstack/localstack-sdk-go` (Go 1.25)
- `generate.go`: `//go:generate` directives for `oapi-codegen`
- `client/`: porcelain entrypoint with ergonomic helpers
  - `config.go`: shared configuration (base URL, HTTP client)
- `client.go`: convenience re-export of the porcelain entrypoint
- `internal/generated/`: generated types and client (do not edit)
- `specs/new-3.0.json`: OpenAPI spec source used for generation (current)

## Usage
1. Install generator (recorded as a tool dependency in `go.mod`, currently v2.5.1):
   ```bash
   go tool oapi-codegen -h   # after go get -tool ...
   ```
2. Generate code:
   ```bash
   go generate ./...
   ```
3. Build or use the stubbed client:
   ```go
   import "github.com/localstack/localstack-sdk-go/client"

   c, err := client.New()
   if err != nil { /* handle */ }
   _ = c
   ```

## Notes
- Generated files are ignored by default (`internal/generated/`).
- Add higher-level helpers in `client/` without touching generated code.
