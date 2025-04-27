# ────────────────────── Stage 1: Builder ──────────────────────
FROM golang:1.21-alpine3.18 AS builder

# Install git so `go mod download` can fetch modules
RUN apk update && apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy remaining source code and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-service ./cmd/server

# ──────────────────── Stage 2: Minimal Runtime ────────────────────
FROM scratch

# Copy the compiled binary from the builder stage
COPY --from=builder /app/webhook-service /webhook-service

# Run the service
ENTRYPOINT ["/webhook-service"]

