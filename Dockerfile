FROM golang:1.25.0@sha256:9e56f0d0f043a68bb8c47c819e47dc29f6e8f5129b8885bed9d43f058f7f3ed6 AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.25.0@sha256:9e56f0d0f043a68bb8c47c819e47dc29f6e8f5129b8885bed9d43f058f7f3ed6 AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:4f6e739881403e7d50f52a4e574c4e3c88266031fd555303ee2f1ba262523d6a AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
