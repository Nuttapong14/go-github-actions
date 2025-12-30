# Backend Dockerfile - Multi-stage build for Go application
# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ ./

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-s -w" to strip debug info and reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/server .

# Runtime stage
FROM alpine:3.21

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Copy migration files
COPY backend/database/migrations ./database/migrations

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health/live || exit 1

# Run the application
CMD ["./server"]
