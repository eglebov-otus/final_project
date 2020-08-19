# syntax = docker/dockerfile:1-experimental

FROM golang:1.14-alpine AS base
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go mod download

FROM golangci/golangci-lint:latest-alpine AS lint-base

FROM base AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN golangci-lint run --timeout 10m0s ./...

FROM base AS unit-test
RUN go test -v ./...

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/app cmd/app/main.go

FROM scratch AS server
COPY --from=build /out/app .
COPY --from=build /src/configs/config.yml configs/
COPY --from=build /src/cache cache
CMD ["./app"]