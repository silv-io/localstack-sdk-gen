package client

import (
	"context"
	"fmt"
	"net/http"

	internalapi "github.com/localstack/localstack-sdk-go/internal/generated/internalapi"
)

type SignalAction string

const (
	SignalRestart SignalAction = "restart"
	SignalKill    SignalAction = "kill"
)

type InternalClient struct {
	gen *internalapi.ClientWithResponses
}

func (c *InternalClient) Restart(ctx context.Context) error {
	return c.SendSignal(ctx, SignalRestart)
}

func (c *InternalClient) Kill(ctx context.Context) error {
	return c.SendSignal(ctx, SignalKill)
}

// SendSignal restarts or terminates the LocalStack session.
func (c *InternalClient) SendSignal(ctx context.Context, action SignalAction) error {
	var signal internalapi.PostLocalstackHealthJSONBodyAction
	switch action {
	case SignalRestart:
		signal = internalapi.Restart
	case SignalKill:
		signal = internalapi.Kill
	default:
		return fmt.Errorf("send signal: unknown action %q", action)
	}

	body := internalapi.PostLocalstackHealthJSONRequestBody{
		Action: &signal,
	}

	resp, err := c.gen.PostLocalstackHealthWithResponse(ctx, body)
	if err != nil {
		return fmt.Errorf("send signal: %w", err)
	}
	return expectStatus(resp.HTTPResponse, http.StatusOK, "send signal", resp.Body)
}

// UpdateHealthInfo updates health information.
func (c *InternalClient) UpdateHealthInfo(ctx context.Context) error {
	resp, err := c.gen.PutLocalstackHealthWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("update health info: %w", err)
	}
	return expectStatus(resp.HTTPResponse, http.StatusOK, "update health info", resp.Body)
}
