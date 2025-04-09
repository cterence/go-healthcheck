FROM golang:1.24.2@sha256:227d106dca555769db9977f33e5d3d27422c5e75af1afc080b92f390c326de80 AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.24.2@sha256:227d106dca555769db9977f33e5d3d27422c5e75af1afc080b92f390c326de80 AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:27769871031f67460f1545a52dfacead6d18a9f197db77110cfc649ca2a91f44 AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
