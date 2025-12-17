package client

import (
	"context"
	"net/http"
)

func makeRequestEditor(cfg Config) func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		if cfg.UserAgent != "" {
			req.Header.Set("User-Agent", cfg.UserAgent)
		}
		if cfg.AuthHeader != "" {
			req.Header.Set("Authorization", cfg.AuthHeader)
		}
		if cfg.APIKey != "" {
			req.Header.Set("x-api-key", cfg.APIKey)
		}

		for k, vs := range cfg.Headers {
			req.Header.Del(k)
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}

		for _, edit := range cfg.RequestEditors {
			if err := edit(ctx, req); err != nil {
				return err
			}
		}
		return nil
	}
}

