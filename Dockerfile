# ───────────── Stage 1: Builder ─────────────────
FROM golang:1.21-slim AS builder

# Install git (for go mod download) and ca-certs
RUN apt-get update \
 && apt-get install -y --no-install-recommends git ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy module definitions, download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code & compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-service ./cmd/server

# ───────────── Stage 2: Minimal Runtime ──────────
FROM gcr.io/distroless/static:nonroot

# Copy the compiled binary
COPY --from=builder /app/webhook-service /webhook-service

# Port & entrypoint
EXPOSE 8080
ENTRYPOINT ["/webhook-service"]


