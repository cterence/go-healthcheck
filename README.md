# go-healthcheck

Simple HTTP health check application in Go.

## Installation

```bash
go install github.com/cterence/go-healthcheck@latest
```

## Usage

Run the app with `./go-healthcheck` and configure via `config.yaml` (example at [config.example.yaml](./config.example.yaml)). Check health at `http://localhost:3000`.

Port can be changed by setting the `GOHC_PORT` env variable.
