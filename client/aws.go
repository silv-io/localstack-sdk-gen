package client

import (
	awsapi "github.com/localstack/localstack-sdk-go/internal/generated/awsapi"
)

// AWSClient wraps the generated AWS proxy client.
type AWSClient struct {
	gen *awsapi.ClientWithResponses
}
