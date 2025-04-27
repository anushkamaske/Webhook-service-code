# ─── Stage 1: Builder ──────────────────────────────────────────────────────────
FROM golang:1.21-bullseye AS builder

# Create and switch to our app directory
WORKDIR /app

# Only copy go.mod and go.sum first, so we can cache deps
COPY go.mod go.sum ./
RUN go mod download

# Now bring in the rest of the source and compile
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags="-s -w" \
             -o webhook-service \
             ./cmd/server

# ─── Stage 2: Minimal Runtime ─────────────────────────────────────────────────
FROM scratch

# Copy the statically-built binary
COPY --from=builder /app/webhook-service /webhook-service

# If you expose a port, document it
EXPOSE 8080

# Run as the entrypoint
ENTRYPOINT ["/webhook-service"]
