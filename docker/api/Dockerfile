# ─── Stage 1: Build ───────────────────────────
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependency files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# ─── Stage 2: Run ─────────────────────────────
FROM alpine:latest

# Needed for OAuth HTTPS calls
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
