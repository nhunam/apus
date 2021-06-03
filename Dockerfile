#####################################
#   STEP 1 build executable binary  #
#####################################
FROM golang:1.15-alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go get -u github.com/swaggo/swag/cmd/swag

COPY . .

ENV ENVIRONMENT=PROD

# Build the binary.
RUN wire;swag init;go build -o main

EXPOSE 8080

# Run the hello binary.
ENTRYPOINT /app/main