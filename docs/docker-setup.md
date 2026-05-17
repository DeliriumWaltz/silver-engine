# Docker Setup for Chat MVP

Use Docker to build and run the chat server locally.

## Quick Start

```bash
# Build and run
docker compose up --build

# Server starts on http://localhost:8080
```

## Commands

```bash
# Build only
docker compose build

# Run in background
docker compose up -d

# Tail logs
docker compose logs -f

# Stop
docker compose down
```

## Image

The `Dockerfile` uses a multi-stage build:
1. **Builder** — `golang:1.22-alpine` compiles the binary
2. **Runtime** — minimal `alpine:3.19` with a non-root user

For production, tag and push the image:
```bash
docker build -t silver-engine/chat-server:latest .
docker run -p 8080:8080 silver-engine/chat-server:latest
```
