# ────────────────────── Stage 1: Builder ──────────────────────
FROM golang:1.21-alpine3.18 AS builder

# 1.1 Install git so `go mod download` can fetch modules
RUN apk update && apk add --no-cache git

WORKDIR /app

# 2. Copy go.mod/go.sum and download deps
COPY go.mod go.sum ./
RUN go mod download

# 3. Copy the rest of the code and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-service ./cmd/server

# ──────────────────── Stage 2: Minimal Runtime ────────────────────
FROM scratch

COPY --from=builder /app/webhook-service /webhook-service

ENTRYPOINT ["/webhook-service"]
