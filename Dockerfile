# ─── Stage 1: Builder ─────────────────────────────────────────────
FROM golang:1.21 AS builder

# Install git & CA certificates so `go mod download` works reliably
RUN apt-get update \
 && apt-get install -y --no-install-recommends git ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-service ./cmd/server

# ─── Stage 2: Minimal Runtime ────────────────────────────────────
FROM gcr.io/distroless/static:nonroot

# Copy CA certs for HTTPS calls (if your service ever calls out)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary
COPY --from=builder /app/webhook-service /webhook-service

EXPOSE 8080
ENTRYPOINT ["/webhook-service"]
