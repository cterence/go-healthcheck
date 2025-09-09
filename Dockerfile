FROM golang:1.25.1@sha256:b773c94fd926723a415b28e615016a49793764e270566a82c6a47afb4c52f7de AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.25.1@sha256:b773c94fd926723a415b28e615016a49793764e270566a82c6a47afb4c52f7de AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:fa15492938650e1a5b87e34d47dc7d99a2b4e8aefd81b931b3f3eb6bb4c1d2f6 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
