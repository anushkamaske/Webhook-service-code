# 1) Builder stage
FROM golang:1.21-alpine AS builder

# Install git so `go mod download` can pull private or public modules
RUN apk add --no-cache git

WORKDIR /app

# Copy module definitions first, download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-service ./cmd/server

# 2) Final minimal image
FROM scratch
COPY --from=builder /app/webhook-service /webhook-service
ENTRYPOINT ["/webhook-service"]
