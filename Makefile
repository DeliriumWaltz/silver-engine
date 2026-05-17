.PHONY: all build test test-unit test-integration test-all clean run docker-build docker-run

BINARY=chat-server
BINARY_PATH=./bin/$(BINARY)

all: build test

build:
	go build -o $(BINARY_PATH) ./cmd/chat-server

run: build
	$(BINARY_PATH)

test:
	go test -v -count=1 ./...

test-unit:
	go test -v -count=1 -short ./...

test-integration:
	go test -v -count=1 -run Integration ./...

test-all: test

clean:
	rm -rf bin/
	go clean -cache

# Docker targets
docker-build:
	docker build -t silver-engine/chat-server:latest .

docker-run: docker-build
	docker run -p 8080:8080 silver-engine/chat-server:latest

docker-compose-up:
	docker compose up --build

docker-compose-down:
	docker compose down
