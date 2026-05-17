# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /build/bin/chat-server ./cmd/chat-server

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user
RUN adduser -D -h /app chat
USER chat
WORKDIR /app

COPY --from=builder /build/bin/chat-server /app/chat-server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/rooms || exit 1

ENTRYPOINT ["/app/chat-server"]
