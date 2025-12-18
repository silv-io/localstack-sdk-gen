package client

import (
	chaosapi "github.com/localstack/localstack-sdk-go/internal/generated/chaosapi"
)

// ChaosClient wraps the generated Chaos client.
type ChaosClient struct {
	gen *chaosapi.ClientWithResponses
}
