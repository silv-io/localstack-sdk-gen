package localstack

// Info APIs
//go:generate go tool oapi-codegen --config config/oapi.info.yaml specs/minimal-3.0.json
//go:generate go tool oapi-codegen --config config/oapi.info.client.yaml specs/minimal-3.0.json

// Internal APIs
//go:generate go tool oapi-codegen --config config/oapi.internal.yaml specs/minimal-3.0.json
//go:generate go tool oapi-codegen --config config/oapi.internal.client.yaml specs/minimal-3.0.json
