# Application configuration
APP_NAME = rePacks
BIN_DIR = bin
UI_DIR = ui
CONFIG_DIR = config
PORT = 8086

# Docker configuration
DOCKER_IMAGE = ${APP_NAME}
DOCKER_TAG = latest
DOCKER_COMPOSE = docker-compose.yml

.PHONY: all build run test clean docker-build docker-run docker-stop ui

all: build

# Build the application
build:
	@echo "Building ${APP_NAME}..."
	@mkdir -p ${BIN_DIR}
	@go build -o ${BIN_DIR}/${APP_NAME} ./cmd/server
	@echo "Build complete: ${BIN_DIR}/${APP_NAME}"

# Run the application
run: build
	@echo "Starting ${APP_NAME} on port ${PORT}..."
	@${BIN_DIR}/${APP_NAME}

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ${BIN_DIR}
	@go clean
	@echo "Clean complete"

# Docker operations
docker-build:
	@echo "Building Docker image ${DOCKER_IMAGE}:${DOCKER_TAG}..."
	@docker compose -f ${DOCKER_COMPOSE} build

docker-run: docker-build
	@echo "Starting containers..."
	@docker compose -f ${DOCKER_COMPOSE} up -d

docker-stop:
	@echo "Stopping containers..."
	@docker compose -f ${DOCKER_COMPOSE} down

docker-logs:
	@docker compose -f ${DOCKER_COMPOSE} logs -f

# UI development
ui:
	@echo "Serving UI on http://localhost:8086"
	@cd ${UI_DIR} && python3 -m http.server 8086

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download

# Format code
fmt:
	@go fmt ./...

# Lint code
lint:
	@golangci-lint run

# Help
help:
	@echo "Available targets:"
	@echo "  all         - Build the application (default)"
	@echo "  build       - Build the application"
	@echo "  run         - Build and run the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Build and run Docker containers"
	@echo "  docker-stop - Stop Docker containers"
	@echo "  docker-logs - View container logs"
	@echo "  ui          - Serve UI for development"
	@echo "  deps        - Install dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  help        - Show this help"