FROM golang:1.25.0@sha256:e5ab8640ed89da51a95d687efdcef75bf1bf7f66368af2818ff7fd34b0bc6539 AS fetch-stage

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download


FROM golang:1.25.0@sha256:e5ab8640ed89da51a95d687efdcef75bf1bf7f66368af2818ff7fd34b0bc6539 AS build-stage
COPY --from=fetch-stage /src /src
COPY main.go pkg/ /src/
COPY pkg/ /src/pkg
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/app


FROM gcr.io/distroless/base-debian12@sha256:4f6e739881403e7d50f52a4e574c4e3c88266031fd555303ee2f1ba262523d6a AS deploy-stage
WORKDIR /app
COPY --from=build-stage /src/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
