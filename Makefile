# Makefile for voidling Discord bot

# Variables
BINARY_NAME=voidling
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_UNIX=$(BINARY_NAME)
CMD_PATH=./cmd/voidling
MIGRATIONS_DIR=./migrations
BUILD_DIR=./build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-trimpath

.PHONY: all build build-windows build-linux build-darwin clean test coverage run install-tools sqlc-generate migrate-up migrate-down migrate-status fmt vet lint help

# Default target
all: clean fmt vet build

## help: Display this help message
help:
	@echo "voidling - Discord Bot Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## build: Build the application for current OS
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BINARY_NAME)"

## build-windows: Build for Windows (amd64)
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_WINDOWS) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_WINDOWS)"

## build-linux: Build for Linux (amd64)
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_UNIX)"

## build-darwin: Build for macOS (amd64 and arm64)
build-darwin:
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX)-darwin-amd64 $(CMD_PATH)
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX)-darwin-arm64 $(CMD_PATH)
	@echo "Build complete for macOS"

## build-all: Build for all platforms
build-all: build-windows build-linux build-darwin
	@echo "All platform builds complete"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(BINARY_WINDOWS) $(BINARY_UNIX)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## run: Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

## dev: Run with live reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not found. Install with: go install github.com/air-verse/air@latest"; \
		exit 1; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

## lint: Run golangci-lint
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running golangci-lint..."; \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install from: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

## tidy: Tidy go modules
tidy:
	@echo "Tidying go modules..."
	$(GOMOD) tidy

## download: Download go modules
download:
	@echo "Downloading dependencies..."
	$(GOMOD) download

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@echo "Installing sqlc..."
	$(GOGET) github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Installing goose..."
	$(GOGET) github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Installing air (live reload)..."
	$(GOGET) github.com/air-verse/air@latest
	@echo "Tools installed successfully"

## sqlc-generate: Generate sqlc code from queries
sqlc-generate:
	@echo "Generating sqlc code..."
	@if command -v sqlc > /dev/null; then \
		sqlc generate; \
		echo "sqlc generation complete"; \
	else \
		echo "sqlc not found. Run 'make install-tools' first"; \
		exit 1; \
	fi

## migrate-up: Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	@if command -v goose > /dev/null; then \
		goose -dir $(MIGRATIONS_DIR) sqlite3 $${DATABASE_PATH:-~/.voidling/voidling.db} up; \
	else \
		echo "goose not found. Run 'make install-tools' first"; \
		exit 1; \
	fi

## migrate-down: Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@if command -v goose > /dev/null; then \
		goose -dir $(MIGRATIONS_DIR) sqlite3 $${DATABASE_PATH:-~/.voidling/voidling.db} down; \
	else \
		echo "goose not found. Run 'make install-tools' first"; \
		exit 1; \
	fi

## migrate-status: Show migration status
migrate-status:
	@echo "Migration status..."
	@if command -v goose > /dev/null; then \
		goose -dir $(MIGRATIONS_DIR) sqlite3 $${DATABASE_PATH:-~/.voidling/voidling.db} status; \
	else \
		echo "goose not found. Run 'make install-tools' first"; \
		exit 1; \
	fi

## migrate-create: Create a new migration (usage: make migrate-create NAME=migration_name)
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@if command -v goose > /dev/null; then \
		goose -dir $(MIGRATIONS_DIR) create $(NAME) sql; \
	else \
		echo "goose not found. Run 'make install-tools' first"; \
		exit 1; \
	fi

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t voidling:latest .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it \
		-v $${PWD}/.env:/app/.env \
		-v $${PWD}/data:/app/data \
		voidling:latest

## init: Initialize project (install tools, download deps, generate code)
init: install-tools download sqlc-generate
	@echo "Project initialized successfully"

## check: Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "All checks passed"

## release: Build optimized release binaries for all platforms
release: clean
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)
	@$(MAKE) build-all
	@echo "Release builds complete in $(BUILD_DIR)/"
