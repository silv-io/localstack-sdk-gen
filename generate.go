package localstack

// Info APIs
//go:generate go tool oapi-codegen --config config/oapi.info.yaml specs/new-3.0.json
//go:generate go tool oapi-codegen --config config/oapi.info.client.yaml specs/new-3.0.json

// Internal APIs
//go:generate go tool oapi-codegen --config config/oapi.internal.yaml specs/new-3.0.json
//go:generate go tool oapi-codegen --config config/oapi.internal.client.yaml specs/new-3.0.json

// Chaos APIs
//go:generate go tool oapi-codegen --config config/oapi.chaos.yaml specs/new-3.0.json
//go:generate go tool oapi-codegen --config config/oapi.chaos.client.yaml specs/new-3.0.json

// AWS service proxy APIs
//go:generate go tool oapi-codegen --config config/oapi.aws.yaml specs/new-3.0.json
//go:generate go tool oapi-codegen --config config/oapi.aws.client.yaml specs/new-3.0.json
