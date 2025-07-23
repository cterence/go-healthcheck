FROM golang:1.24.5@sha256:a3bb6cd5f068b34961d60dcd7fc51fb70df7b74b4aa89ac480fc38a6ccba265e AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.24.5@sha256:a3bb6cd5f068b34961d60dcd7fc51fb70df7b74b4aa89ac480fc38a6ccba265e AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:201ef9125ff3f55fda8e0697eff0b3ce9078366503ef066653635a3ac3ed9c26 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
