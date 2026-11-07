.PHONY: all build build-web test bench clean run

BINARY_NAME=ogcraft
VERSION=1.0.0

all: test build

build-web:
	@echo "Building React 19 Visual Designer..."
	cd web && npm install && npm run build

build:
	@echo "Building single static binary..."
	go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BINARY_NAME) ./cmd/ogcraft

test:
	@echo "Running tests..."
	go test -v -race ./...

bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./internal/renderer/... ./internal/cache/...

run:
	go run ./cmd/ogcraft serve --port 8080

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -rf data/cache/*
