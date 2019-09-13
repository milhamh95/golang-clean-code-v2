# ## Builder
# FROM golang:1.13.0-alpine3.10 as builder

# # Install package
# RUN apk update && apk upgrade \
#     && apk --update --no-cache add git make \
#     && rm -f /var/cache/apk/*

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]

# ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o employee" -command="./employee http"

FROM golang:1.12-alpine
RUN apk add --update --no-cache git

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o employee" -command="./employee http"
