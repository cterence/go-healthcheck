FROM golang:1.25.5@sha256:97be07314ef2af5f56d22c3bb608c4cffa2a92b3c8252e9f674081ed8217f75b AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.25.5@sha256:97be07314ef2af5f56d22c3bb608c4cffa2a92b3c8252e9f674081ed8217f75b AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:f5a3067027c2b322cd71b844f3d84ad3deada45ceb8a30f301260a602455070e AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
