## Builder
FROM golang:1.13.0-alpine3.10 as builder

# Install package
RUN apk update && apk upgrade \
    && apk --update --nocache add git gcc make tzdata \
    && rm -f /var/cache/apk/*

WORKDIR /app

COPY go.mod go.sum ./
