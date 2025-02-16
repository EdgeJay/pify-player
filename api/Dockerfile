# Stage 1: Build the Go application
FROM golang:1.24.0 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies (caching step)
RUN go mod download

# Build the Go application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o pify-player-api .

# Stage 2: Create a lightweight runtime container
FROM alpine:latest

# Set environment variables
ENV PORT=5000

# Set the working directory in the final container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/pify-player-api .

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Expose the application port
EXPOSE 5000

# Run the Go application
ENTRYPOINT ["/app/pify-player-api"]
