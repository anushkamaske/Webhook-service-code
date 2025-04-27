# ──────────────── Stage 1: Builder ────────────────
FROM golang:1.21-alpine AS builder

# 1) Install git so go mod can fetch modules
RUN apk add --no-cache git

WORKDIR /app

# 2) Copy only go.mod and go.sum, then download deps
COPY go.mod go.sum ./
RUN go mod download

# 3) Copy the rest of your code and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-service ./cmd/server

# ────────────── Stage 2: Minimal Runtime ──────────────
FROM scratch

# Copy the statically‐built binary
COPY --from=builder /app/webhook-service /webhook-service

# Run it!
ENTRYPOINT ["/webhook-service"]
