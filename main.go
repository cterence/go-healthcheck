package main

import (
	"encoding/json"
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
	healthPostgres "github.com/hellofresh/health-go/v5/checks/postgres"
	healthRedis "github.com/hellofresh/health-go/v5/checks/redis"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Name           string   `yaml:"name"`
	Version        string   `yaml:"version"`
	Timeout        int      `yaml:"timeout"`
	URLs           []string `yaml:"urls"`
	PostgreSQLURIs []string `yaml:"postgresqlURIs"`
	RedisURIs      []string `yaml:"redisURIs"`
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
		url, err := url.ParseRequestURI(u)
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
			log.Fatalf("failed to register http health check %s: %v", url, err)
		}
	}

	for _, u := range config.PostgreSQLURIs {
		url, err := url.ParseRequestURI(u)
		if err != nil {
			log.Fatalf("failed to parse url %s: %v", u, err)
		}

		err = h.Register(health.Config{
			Name:      url.Host,
			Timeout:   time.Second * time.Duration(config.Timeout),
			SkipOnErr: false,
			Check: healthPostgres.New(healthPostgres.Config{
				DSN: u,
			}),
		})
		if err != nil {
			log.Fatalf("failed to register postgresql health check %s: %v", url, err)
		}
	}

	for _, u := range config.RedisURIs {
		url, err := url.ParseRequestURI(u)
		if err != nil {
			log.Fatalf("failed to parse url %s: %v", u, err)
		}

		err = h.Register(health.Config{
			Name:      url.Host,
			Timeout:   time.Second * time.Duration(config.Timeout),
			SkipOnErr: false,
			Check: healthRedis.New(healthRedis.Config{
				DSN: u,
			}),
		})
		if err != nil {
			log.Fatalf("failed to register redis health check %s: %v", url, err)
		}
	}

	logger := httplog.NewLogger("go-healthcheck", httplog.Options{
		LogLevel:        slog.LevelInfo,
		Concise:         true,
		RequestHeaders:  true,
		TimeFieldFormat: time.RFC3339,
	})

	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(time.Second * 10))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		c := h.Measure(r.Context())

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		code := http.StatusOK
		if c.Status == "Unavailable" {
			code = http.StatusServiceUnavailable
		}
		logger.Logger.Info(string(data))
		w.WriteHeader(code)
		w.Write(data)
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
