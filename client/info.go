package client

import (
	"context"
	"fmt"
	"net/http"

	infoapi "github.com/localstack/localstack-sdk-go/internal/generated/infoapi"
)

type InfoClient struct {
	gen *infoapi.ClientWithResponses
}

func (c *InfoClient) Health(ctx context.Context) (*Health, error) {
	resp, err := c.gen.GetLocalstackHealthWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("get health: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get health", resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get health: empty body")
	}
	return (*Health)(resp.JSON200), nil
}

func (c *InfoClient) SessionInfo(ctx context.Context) (*Info, error) {
	resp, err := c.gen.GetLocalstackInfoWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("get session info: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get session info", resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get session info: empty body")
	}
	return (*Info)(resp.JSON200), nil
}

func (c *InfoClient) LicenseInfo(ctx context.Context) (*LicenseInfo, error) {
	resp, err := c.gen.GetLocalstackLicenseinfoWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("get license info: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get license info", resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get license info: empty body")
	}
	return (*LicenseInfo)(resp.JSON200), nil
}

func (c *InfoClient) Diagnose(ctx context.Context) (*Diagnose, error) {
	resp, err := c.gen.GetLocalstackDiagnoseWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("get diagnose: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get diagnose", resp.Body); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get diagnose: empty body")
	}
	return (*Diagnose)(resp.JSON200), nil
}

// Usage returns raw usage statistics bytes.
func (c *InfoClient) Usage(ctx context.Context) ([]byte, error) {
	resp, err := c.gen.GetLocalstackUsageWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("get usage: %w", err)
	}
	if err := expectStatus(resp.HTTPResponse, http.StatusOK, "get usage", resp.Body); err != nil {
		return nil, err
	}
	return resp.Body, nil
}
