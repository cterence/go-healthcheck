# go-healthcheck

Simple health check application based on the [hellofresh/health-go](https://github.com/hellofresh/health-go) library.

## Supported targets

- HTTP(S)
- PostgreSQL
- Redis

## Installation

Install binary with go:

```bash
go install github.com/cterence/go-healthcheck@latest
```

Or run with Docker:

```bash
docker run -v ./config.yaml:/app/config.yaml ghcr.io/cterence/go-healthcheck:latest
```

## Usage

```bash
# create a config file
mv config.example.yaml config.yaml
# run example docker compose
docker compose up -d
# query go-healthcheck
curl localhost:3000
# {"status":"OK","timestamp":"2025-03-19T21:25:35.401294936Z","component":{"name":"mychecks","version":"1.0"}}
```

Listening port can be changed by setting the `GOHC_PORT` env variable.
