//go:build integration

package integration_test

import (
	"context"
	"fmt"

	"github.com/localstack/localstack-sdk-go/client"
	"github.com/localstack/localstack-sdk-go/internal/testsetup"
)

// Example demonstrating how to start LocalStack via docker and call Info APIs.
func Example_integration() {
	ctx := context.Background()
	inst, err := testsetup.StartLocalStackDocker(ctx)
	if err != nil {
		fmt.Println("start localstack:", err)
		return
	}
	defer inst.Stop(context.Background())

	c, err := client.New(client.WithBaseURL(inst.Endpoint))
	if err != nil {
		fmt.Println("new client:", err)
		return
	}

	if _, err := c.Info.Health(ctx); err != nil {
		fmt.Println("health:", err)
		return
	}

	info, err := c.Info.SessionInfo(ctx)
	if err != nil {
		fmt.Println("session info:", err)
		return
	}

	fmt.Println("session info retrieved, edition:", info.Edition)
	// Output:
	// session info retrieved, edition: community
}
