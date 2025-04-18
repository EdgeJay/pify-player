# Stage 1: Build the Go application
FROM golang:1.24.0 AS builder

ARG BUNDEBUG
ARG PORT
ARG CORS_ORIGINS
ARG DB_FILE
ARG SPOTIFY_CLIENT_ID
ARG SPOTIFY_CLIENT_SECRET
ARG SPOTIFY_REDIRECT_URI
ARG CALLBACK_DEST
ARG ALLOW_SHELL_COMMANDS

ENV BUNDEBUG=${BUNDEBUG}
ENV PORT=${PORT}
ENV CORS_ORIGINS=${CORS_ORIGINS}
ENV DB_FILE=${DB_FILE}
ENV SPOTIFY_CLIENT_ID=${SPOTIFY_CLIENT_ID}
ENV SPOTIFY_CLIENT_SECRET=${SPOTIFY_CLIENT_SECRET}
ENV SPOTIFY_REDIRECT_URI=${SPOTIFY_REDIRECT_URI}
ENV CALLBACK_DEST=${CALLBACK_DEST}
ENV ALLOW_SHELL_COMMANDS=${ALLOW_SHELL_COMMANDS}

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies (caching step)
RUN go mod download

# Build the Go application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o pify-player-api ./cmd/api/main.go

# Stage 2: Create a lightweight runtime container
FROM alpine:latest

# Set the working directory in the final container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/pify-player-api .
COPY --from=builder /app/certs/ ./certs

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Run the Go application
ENTRYPOINT ["/app/pify-player-api"]
