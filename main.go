package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
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

	logger := httplog.NewLogger("go-healthcheck", httplog.Options{
		LogLevel:        slog.LevelInfo,
		Concise:         true,
		RequestHeaders:  true,
		TimeFieldFormat: time.RFC3339,
	})

	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(time.Second * 10))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.HandlerFunc(w, r)
	})

	port := os.Getenv("GOHC_PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("Listening on port %s\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}
}
