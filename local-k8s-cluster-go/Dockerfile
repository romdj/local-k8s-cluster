# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o k3s-manager ./cmd

# Final stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ca-certificates kubectl

# Create non-root user
RUN adduser -D -s /bin/sh k3s-manager

# Copy binary from builder
COPY --from=builder /app/k3s-manager /usr/local/bin/k3s-manager

# Set ownership and permissions
RUN chown k3s-manager:k3s-manager /usr/local/bin/k3s-manager && \
    chmod +x /usr/local/bin/k3s-manager

# Switch to non-root user
USER k3s-manager

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/k3s-manager"]