package target

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cterence/go-healthcheck/pkg/config"
	healthhttp "github.com/cterence/go-healthcheck/pkg/lib/http"
	"github.com/hellofresh/health-go/v5"
)

type HTTP struct {
	URL *url.URL
}

func (t *HTTP) New(uri string) error {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return fmt.Errorf("failed to parse url %s: %v", u, err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid scheme for http target: %s", u.Scheme)
	}

	t.URL = u
	return nil
}

func (t *HTTP) Register(h *health.Health, c *config.Config) error {
	httpConfig := healthhttp.Config{
		URL:                          t.URL.String(),
		RequestTimeout:               time.Second * time.Duration(c.Timeout),
		HTTPStatusCodeErrorThreshold: c.HTTPStatusCodeErrorThreshold,
	}

	if c.HTTPClientCertPath != "" && c.HTTPClientKeyPath != "" {
		cert, err := tls.LoadX509KeyPair(c.HTTPClientCertPath, c.HTTPClientKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load client cert and key: %v", err)
		}
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			},
			Timeout: time.Second * time.Duration(c.Timeout),
		}
		httpConfig.Client = client
		fmt.Printf("Using client cert and key from %s and %s\n", c.HTTPClientCertPath, c.HTTPClientKeyPath)
	}

	if err := h.Register(health.Config{
		Name:      t.URL.Host,
		Timeout:   time.Second * time.Duration(c.Timeout),
		SkipOnErr: false,
		Check:     healthhttp.New(httpConfig),
	}); err != nil {
		return fmt.Errorf("failed to register HTTP health check %s: %v", t.URL, err)
	}
	return nil
}

func (t *HTTP) String() string {
	return t.URL.Redacted()
}
