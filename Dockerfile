# Stage 1: Build the binary
FROM golang:1.22.5-alpine AS builder

# Enable Go modules
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary (disable CGO to ensure static binary)
RUN CGO_ENABLED=0 go build -o engine-app ./cmd/app

# Stage 2: Create a minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Copy the binary from builder stage
COPY --from=builder /app/engine-app /engine-app

# Expose the port if needed (optional, for documentation)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/engine-app"]
