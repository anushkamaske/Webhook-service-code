# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Install git so `go mod download` can pull modules
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency definitions and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-service ./cmd/server

# Stage 2: Create minimal runtime image
FROM scratch

# Copy the built binary from the builder stage
COPY --from=builder /app/webhook-service /webhook-service

# Run the webhook-service binary by default
ENTRYPOINT ["/webhook-service"]
