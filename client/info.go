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

// Health checks whether LocalStack is reachable.
func (c *InfoClient) Health(ctx context.Context) error {
	resp, err := c.gen.GetLocalstackHealthWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("get health: %w", err)
	}
	return expectStatus(resp.HTTPResponse, http.StatusOK, "get health", resp.Body)
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
