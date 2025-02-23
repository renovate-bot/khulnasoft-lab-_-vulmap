# Build
FROM golang:1.21-alpine AS build-env
RUN apk add build-base
WORKDIR /app
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build ./cmd/vulmap

# Release
FROM alpine:3.18.4
RUN apk -U upgrade --no-cache \
    && apk add --no-cache bind-tools chromium ca-certificates
COPY --from=build-env /app/vulmap /usr/local/bin/

ENTRYPOINT ["vulmap"]
