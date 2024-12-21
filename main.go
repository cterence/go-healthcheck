package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hellofresh/health-go/v5"
	healthHttp "github.com/hellofresh/health-go/v5/checks/http"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Timeout int      `yaml:"timeout"`
	URLs    []string `yaml:"urls"`
}

func main() {
	configYaml, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}

	config := Config{}

	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML config: %v", err)
	}

	h, _ := health.New(health.WithComponent(health.Component{
		Name:    config.Name,
		Version: config.Version,
	}))

	for _, u := range config.URLs {
		url, err := url.Parse(u)
		if err != nil {
			log.Fatalf("failed to parse url %s: %v", u, err)
		}

		err = h.Register(health.Config{
			Name:      url.Host,
			Timeout:   time.Second * time.Duration(config.Timeout),
			SkipOnErr: false,
			Check: healthHttp.New(healthHttp.Config{
				URL:            url.String(),
				RequestTimeout: time.Second * time.Duration(config.Timeout),
			}),
		})
		if err != nil {
			log.Fatalf("failed to register health check %s: %v", url, err)
		}
	}

	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Header.Add(w.Header(), "Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status": {"server": "OK"}}`))
		if err != nil {
			log.Fatalf("failed to write response for /health endpoint: %v", err)
		}
	}))
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		h.HandlerFunc(w, r)
		took := time.Since(t)
		addr := r.RemoteAddr
		if ip, exists := r.Header["X-Real-Ip"]; exists {
			addr = ip[0]
		}
		slog.Info(fmt.Sprintf("%s - %s - %s", addr, r.UserAgent(), took.Abs().Round(time.Millisecond).String()))
	})
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("failed to listen on port 3000: %v", err)
	}
}
