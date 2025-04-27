# Multi‐stage build for the Go server
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-service ./cmd/server

FROM scratch
COPY --from=builder /app/webhook-service /webhook-service
ENTRYPOINT ["/webhook-service"]
