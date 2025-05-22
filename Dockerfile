FROM golang:1.24.3@sha256:02a22753ab3426d91ba5ba6f4dfb4ac2454f19b05afdb18d61ab02cbf1a2dffe AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.24.3@sha256:02a22753ab3426d91ba5ba6f4dfb4ac2454f19b05afdb18d61ab02cbf1a2dffe AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:cef75d12148305c54ef5769e6511a5ac3c820f39bf5c8a4fbfd5b76b4b8da843 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
