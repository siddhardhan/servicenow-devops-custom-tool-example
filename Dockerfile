# Build stage
FROM golang:1.21-alpine AS builder

# Install required tools and dependencies
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Add required dependencies
RUN go get github.com/swaggo/swag/cmd/swag@v1.8.1 && \
    go get github.com/swaggo/gin-swagger@v1.4.1 && \
    go get github.com/swaggo/files@v0.0.0-20220728132757-551d4a08d97a

# Copy source code
COPY . .

# Generate Swagger documentation
RUN go mod tidy && \
    /go/bin/swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o evidence-service -ldflags="-w -s" -trimpath

# Final stage
FROM alpine:3.18

# Add necessary certificates for HTTPS and timezone data
RUN apk --no-cache add \
    ca-certificates \
    tzdata

WORKDIR /app

# Copy the binary and swagger docs from builder
COPY --from=builder /app/evidence-service .
COPY --from=builder /app/docs/ ./docs/

# Create non-root user
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set environment variables
ENV GIN_MODE=release
ENV TZ=UTC

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/v1/evidences?controlId=1234 || exit 1

# Expose port 8080
EXPOSE 8080

# Run the service
ENTRYPOINT ["./evidence-service"]
