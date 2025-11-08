FROM golang:1.25.4@sha256:6ca9eb0b32a4bd4e8c98a4a2edf2d7c96f3ea6db6eb4fc254eef6c067cf73bb4 AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.25.4@sha256:6ca9eb0b32a4bd4e8c98a4a2edf2d7c96f3ea6db6eb4fc254eef6c067cf73bb4 AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:9e9b50d2048db3741f86a48d939b4e4cc775f5889b3496439343301ff54cdba8 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
