// Original code from https://github.com/hellofresh/health-go/blob/master/checks/http/check.go
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const defaultRequestTimeout = 5 * time.Second

type Config struct {
	URL                          string
	RequestTimeout               time.Duration
	HTTPStatusCodeErrorThreshold int
	Client                       *http.Client
}

func New(config Config) func(ctx context.Context) error {
	if config.RequestTimeout == 0 {
		config.RequestTimeout = defaultRequestTimeout
	}

	client := config.Client
	if client == nil {
		client = http.DefaultClient
	}

	return func(ctx context.Context) error {
		req, err := http.NewRequest(http.MethodGet, config.URL, nil)
		if err != nil {
			return fmt.Errorf("creating the request for the health check failed: %w", err)
		}

		ctx, cancel := context.WithTimeout(ctx, config.RequestTimeout)
		defer cancel()

		req.Header.Set("Connection", "close")
		req = req.WithContext(ctx)

		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("making the request for the health check failed: %w", err)
		}
		defer res.Body.Close() //nolint:errcheck

		fmt.Printf("Response status: %d\n", res.StatusCode)

		threshold := config.HTTPStatusCodeErrorThreshold
		if threshold == 0 {
			threshold = http.StatusInternalServerError
		}

		if res.StatusCode >= threshold {
			return fmt.Errorf("remote service is not available at the moment: %d", res.StatusCode)
		}

		return nil
	}
}
