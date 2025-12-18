package client

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Operation  string
	StatusCode int
	Status     string
	Body       []byte
}

func (e *HTTPError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if len(e.Body) == 0 {
		return fmt.Sprintf("%s: unexpected status %d", e.Operation, e.StatusCode)
	}
	return fmt.Sprintf("%s: unexpected status %d: %s", e.Operation, e.StatusCode, truncateBody(e.Body, 4<<10))
}

func expectStatus(resp *http.Response, want int, operation string, body []byte) error {
	if resp == nil {
		return fmt.Errorf("%s: missing http response", operation)
	}
	if resp.StatusCode == want {
		return nil
	}
	return &HTTPError{
		Operation:  operation,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
	}
}

func truncateBody(body []byte, limit int) string {
	if limit <= 0 || len(body) <= limit {
		return string(body)
	}
	return string(body[:limit]) + "…"
}

