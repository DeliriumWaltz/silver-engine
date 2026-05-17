# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/chat-server ./cmd/chat-server

# Runtime stage — minimal image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/chat-server /chat-server

EXPOSE 8080

ENTRYPOINT ["/chat-server"]
