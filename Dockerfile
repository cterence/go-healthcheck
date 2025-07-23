FROM golang:1.24.5@sha256:267159cb984d1d034fce6e9db8641bf347f80e5f2e913561ea98c40d5051cb67 AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.24.5@sha256:267159cb984d1d034fce6e9db8641bf347f80e5f2e913561ea98c40d5051cb67 AS build-stage
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
