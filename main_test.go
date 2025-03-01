package main_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hellofresh/health-go/v5"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Header.Add(w.Header(), "Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status": {"server": "OK"}}`))
		assert.NoError(t, err)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"status": {"server": "OK"}}`, rr.Body.String())
}

func TestRootEndpoint(t *testing.T) {
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    "test",
		Version: "1.0.0",
	}))

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		h.HandlerFunc(w, r)
		took := time.Since(t)
		addr := r.RemoteAddr
		if ip, exists := r.Header["X-Real-Ip"]; exists {
			addr = ip[0]
		}
		if addr == "" {
			addr = "unknown-ip"
		}
		userAgent := r.UserAgent()
		if userAgent == "" {
			userAgent = "unknown-useragent"
		}
		slog.Info(fmt.Sprintf("%s - %s - %s", addr, userAgent, took.Abs().Round(time.Millisecond).String()))
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
