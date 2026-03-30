# ─── Stage 1: Build ───────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/scheduler

# ─── Stage 2: Run ─────────────────────────────
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/scheduler .

CMD ["./scheduler"]
