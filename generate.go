package localstack

// Shared types (optional)
//go:generate go tool oapi-codegen --config config/oapi.types.yaml specs/localstack-spec-latest.yml

// LocalStack core (admin/misc, excluding pods/chaos/replicator)
//go:generate go tool oapi-codegen --config config/oapi.localstack.yaml specs/localstack-spec-latest.yml
//go:generate go tool oapi-codegen --config config/oapi.localstack.client.yaml specs/localstack-spec-latest.yml

// Chaos APIs
//go:generate go tool oapi-codegen --config config/oapi.chaos.yaml specs/localstack-spec-latest.yml
//go:generate go tool oapi-codegen --config config/oapi.chaos.client.yaml specs/localstack-spec-latest.yml

// Replicator APIs
//go:generate go tool oapi-codegen --config config/oapi.replicator.yaml specs/localstack-spec-latest.yml
//go:generate go tool oapi-codegen --config config/oapi.replicator.client.yaml specs/localstack-spec-latest.yml

// Cloud Pods APIs
//go:generate go tool oapi-codegen --config config/oapi.pods.yaml specs/localstack-spec-latest.yml
//go:generate go tool oapi-codegen --config config/oapi.pods.client.yaml specs/localstack-spec-latest.yml

// AWS service proxy APIs
//go:generate go tool oapi-codegen --config config/oapi.aws.yaml specs/localstack-spec-latest.yml
//go:generate go tool oapi-codegen --config config/oapi.aws.client.yaml specs/localstack-spec-latest.yml
