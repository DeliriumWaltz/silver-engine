# Silver Engine — Multiplayer Chat with AI Coding Agents

A multiplayer chat application with AI agents that can write code and commit it to a central repository.

## MVP — Chat Interface

The initial MVP is a Go-based chat server with:

- **In-memory storage** (pluggable — swap for a DB-backed `Store` implementation later)
- **REST API** for rooms and messages
- **WebSocket stub** ready for real-time messaging
- **Full test coverage** with unit tests + integration tests that spin up the binary

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/rooms` | List all rooms |
| POST | `/api/rooms` | Create a room |
| GET | `/api/rooms/{roomID}` | Get a room |
| GET | `/api/rooms/{roomID}/messages` | Get messages in a room |
| POST | `/api/rooms/{roomID}/messages` | Send a message |
| GET | `/api/rooms/{roomID}/ws` | WebSocket endpoint (stub) |

### Project Structure

```
├── cmd/chat-server/       # Main binary entry point
├── internal/
│   ├── api/               # HTTP handlers & WebSocket hub
│   ├── chat/              # Business logic service
│   ├── models/            # Shared data models
│   └── storage/           # Storage interface + in-memory impl
├── tests/integration/     # Binary-level integration tests
├── docs/                  # Design documents
├── go.mod
├── Makefile
├── Dockerfile             # Multi-stage Docker build
└── docker-compose.yml     # One-command startup
```

### Quick Start

```bash
# Run tests
make test

# Build and run the server
make run
# Server starts on :8080 (or $PORT)

# Or with Docker
docker compose up --build
# Server starts on http://localhost:8080
```

### Docker

```bash
# Build and run with compose
docker compose up --build

# Build image directly
docker build -t silver-engine/chat-server .

# Run standalone
docker run -p 8080:8080 silver-engine/chat-server
```

### Storage Interface

The `storage.Store` interface makes it easy to swap out the in-memory store:

```go
type Store interface {
    CreateRoom(name string) (*models.Room, error)
    GetRoom(id string) (*models.Room, error)
    ListRooms() ([]*models.Room, error)
    SaveMessage(msg *models.Message) error
    GetMessages(roomID string) ([]*models.Message, error)
}
```

Just implement it with PostgreSQL, Redis, or whatever you like!
