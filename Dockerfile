# VecScan Dockerfile
# Build: docker build -t vecscan .
# Run:   docker run --rm --network host vecscan [target]

FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev libpcap-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o vecscan ./cmd/vecscan

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache libpcap

# Copy binary from builder
COPY --from=builder /app/vecscan /usr/local/bin/vecscan

# Make executable
RUN chmod +x /usr/local/bin/vecscan

# Set entrypoint
ENTRYPOINT ["vecscan"]