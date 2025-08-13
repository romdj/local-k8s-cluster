# Build stage
FROM golang:1.24-alpine AS builder

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

# Install runtime dependencies and kubectl
RUN apk add --no-cache ca-certificates curl && \
    ARCH=$(uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/') && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${ARCH}/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl

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