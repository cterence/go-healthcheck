package target

import (
	"fmt"
	"net/url"
	"time"

	"github.com/cterence/go-healthcheck/pkg/config"
	"github.com/hellofresh/health-go/v5"
	"github.com/hellofresh/health-go/v5/checks/http"
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
	if err := h.Register(health.Config{
		Name:      t.URL.Host,
		Timeout:   time.Second * time.Duration(c.Timeout),
		SkipOnErr: false,
		Check: http.New(http.Config{
			URL:            t.URL.String(),
			RequestTimeout: time.Second * time.Duration(c.Timeout),
		}),
	}); err != nil {
		return fmt.Errorf("failed to register HTTP health check %s: %v", t.URL, err)
	}
	return nil
}

func (t *HTTP) String() string {
	return t.URL.Redacted()
}
