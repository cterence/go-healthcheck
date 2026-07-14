FROM golang:1.26.5@sha256:d52df9c279840adf958d017ebb275651ed8338b953d39817bc3633a2e6b1bbcc AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.26.5@sha256:d52df9c279840adf958d017ebb275651ed8338b953d39817bc3633a2e6b1bbcc AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:348dac1808083ccc3366399d6db835875b4eaf7c9b694783f5a3f353c4b58a28 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
