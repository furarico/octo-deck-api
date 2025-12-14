# Build stage
FROM golang:1.25.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install Task
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Copy source code
COPY . .

# Generate code from OpenAPI spec
RUN task generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -tags timetzdata -o ./bin/main ./cmd/server/main.go

# Runtime stage
FROM alpine:3.22.2

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/main /app/main

# Copy OpenAPI spec
COPY --from=builder /app/openapi /app/openapi

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
