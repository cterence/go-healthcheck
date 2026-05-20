FROM golang:1.26.3@sha256:6df14f4a4bc9d979a3721f488981e0d1b318006377e473ed23d026796f5f4c0a AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.26.3@sha256:6df14f4a4bc9d979a3721f488981e0d1b318006377e473ed23d026796f5f4c0a AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:58695f439f772a00009c8f6be4c183f824c1f556d74b313c30900f167e4772f8 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
