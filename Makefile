.PHONY: build up down test test-unit test-integration clean

# Build Docker images for REST API and gRPC store
build:
	docker-compose build

# Start services in detached mode
up:
	docker-compose up -d

# Stop and remove containers, networks
down:
	docker-compose down

# Run all integration tests
test: test-integration

# Run integration tests inside the containers environment
test-integration: up
	@echo "Waiting 5 seconds for services to be ready..."
	sleep 5
	go test ./integration_test/

# Clean build artifacts (optional)
clean:
	docker-compose down --rmi all --volumes --remove-orphans
	go clean ./...
