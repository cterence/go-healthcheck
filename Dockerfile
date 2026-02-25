FROM golang:1.26.0@sha256:a9c4aaca2982a34e9e5b8d5e906f29cf797a893655182bcf8b90d65071a730ea AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.26.0@sha256:a9c4aaca2982a34e9e5b8d5e906f29cf797a893655182bcf8b90d65071a730ea AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:937c7eaaf6f3f2d38a1f8c4aeff326f0c56e4593ea152e9e8f74d976dde52f56 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
