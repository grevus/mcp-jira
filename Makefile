.PHONY: build test test-integration lint clean

build:
	go build -o bin/mcp-issues ./cmd/mcp-issues
	go build -o bin/mcp-issues-index ./cmd/mcp-issues-index

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

lint:
	go vet ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"

clean:
	rm -rf bin/ mcp-issues mcp-issues-index
