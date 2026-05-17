.PHONY: all build test test-unit test-integration test-all clean run

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

# TDD watch mode helper (requires fswatch or similar)
# watch:
#	while inotifywait -r -e modify .; do make test; done