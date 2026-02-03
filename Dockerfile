FROM golang:1.25.6@sha256:0c87ea6991c06552ca5f516e3aeb434056bac3b674f32f612691692668e57074 AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.25.6@sha256:0c87ea6991c06552ca5f516e3aeb434056bac3b674f32f612691692668e57074 AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:347a41e7f263ea7f7aba1735e5e5b1439d9e41a9f09179229f8c13ea98ae94cf AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
