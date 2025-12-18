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

func (c *InternalClient) InitScripts(ctx context.Context) (*InitScripts, error) {
	resp, err := c.gen.GetLocalstackInitWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("get init scripts: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get init scripts", resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get init scripts: empty body")
	}
	return (*InitScripts)(resp.JSON200), nil
}

func (c *InternalClient) InitScriptsStage(ctx context.Context, stage string) (*InitScripts, error) {
	resp, err := c.gen.GetLocalstackInitStageWithResponse(ctx, stage)
	if err != nil {
		return nil, fmt.Errorf("get init scripts stage: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get init scripts stage", resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get init scripts stage: empty body")
	}
	return (*InitScripts)(resp.JSON200), nil
}

func (c *InternalClient) Plugins(ctx context.Context) (*Plugins, error) {
	resp, err := c.gen.GetLocalstackPluginsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("get plugins: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get plugins", resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get plugins: empty body")
	}
	return (*Plugins)(resp.JSON200), nil
}

func (c *InternalClient) Certificates(ctx context.Context) (*CertificateList, error) {
	resp, err := c.gen.GetListCertificatesWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("get certificates: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get certificates", resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get certificates: empty body")
	}
	return (*CertificateList)(resp.JSON200), nil
}
