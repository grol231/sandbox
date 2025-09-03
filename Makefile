.PHONY: build test clean run docker-build docker-run

# Build the application
build:
	go build -o build/worker ./cmd/worker

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f build/worker
	rm -f coverage.out coverage.html

# Run the application locally
run:
	go run ./cmd/worker -config=configs/config.yaml

# Download dependencies
deps:
	go mod download
	go mod verify

# Build Docker image
docker-build: build
	docker build -f build/Dockerfile -t starline-worker:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 starline-worker:latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Vet code
vet:
	go vet ./...