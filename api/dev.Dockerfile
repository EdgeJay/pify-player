# Stage 1: Build the Go application
FROM golang:1.24.0 AS builder

# Set the working directory inside the container
WORKDIR /app

RUN go install github.com/air-verse/air@latest

# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
