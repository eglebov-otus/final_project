FROM golang:1.14 AS base
WORKDIR /src
COPY . .
RUN go mod download

FROM golangci/golangci-lint:latest-alpine AS lint-base

FROM base AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN golangci-lint run ./...

FROM base AS unit-test
RUN go test -v -race -count 100 ./...

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/app cmd/app/main.go

FROM alpine:latest AS server
COPY --from=build /out/app .
COPY --from=build /src/configs/config.yml configs/
COPY --from=build /src/cache cache
CMD ["./app"]