package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cterence/go-healthcheck/pkg/config"
	"github.com/cterence/go-healthcheck/pkg/router"
	"github.com/cterence/go-healthcheck/pkg/target"
	"github.com/hellofresh/health-go/v5"
)

func main() {
	config := &config.Config{}
	err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	h, err := health.New(health.WithComponent(health.Component{
		Name:    config.Name,
		Version: config.Version,
	}))
	if err != nil {
		log.Fatalf("failed to create health check container: %v", err)
	}

	for _, endpoint := range config.Targets.HTTP {
		t := &target.HTTP{}
		if err := target.Register(t, endpoint, h, config); err != nil {
			log.Fatalf("failed to register http target: %v", err)
		}
	}

	for _, endpoint := range config.Targets.PostgreSQL {
		t := &target.PostgreSQL{}
		if err := target.Register(t, endpoint, h, config); err != nil {
			log.Fatalf("failed to register postgresql target: %v", err)
		}
	}

	for _, endpoint := range config.Targets.Redis {
		t := &target.Redis{}
		if err := target.Register(t, endpoint, h, config); err != nil {
			log.Fatalf("failed to register redis target: %v", err)
		}
	}

	r := router.New(h)

	fmt.Printf("Listening on port %s\n", config.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", config.Port, err)
	}
}
