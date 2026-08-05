FROM golang:1.26.5@sha256:6cd10a6fcc5eadd62008fc2ad8056b38971cafd42f44d55297f18be8adc86410 AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.26.5@sha256:6cd10a6fcc5eadd62008fc2ad8056b38971cafd42f44d55297f18be8adc86410 AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:62730825d3cf03571e0a1b8f014748de94d0404500f063593b614c23da38841d AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
