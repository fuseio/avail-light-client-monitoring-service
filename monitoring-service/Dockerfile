# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/server ./cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite

# Create a non-root user
RUN adduser -D -g '' appuser

# Create directory for SQLite database
RUN mkdir -p /data && chown -R appuser:appuser /data

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/server .

# Use non-root user
USER appuser

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./server"]