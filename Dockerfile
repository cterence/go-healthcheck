FROM golang:1.27.0@sha256:fb7313a7349714c1df7476aae7c8c271132fc30bfecd8bb58678b4f9fb311aff AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.27.0@sha256:fb7313a7349714c1df7476aae7c8c271132fc30bfecd8bb58678b4f9fb311aff AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:76b3162a31477bca4a245b836c624f4c4a1a3705e99b9003907d992bec2c4bca AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
