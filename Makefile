.PHONY: help test unit-test integration-test coverage lint fmt run build clean docker-up docker-down dashboard dashboard-deps openapi-validate

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: unit-test integration-test ## Run all tests

unit-test: ## Run unit tests
	@echo "Running unit tests..."
	go test -v -race -cover ./internal/domain/...

integration-test: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -race -timeout 5m ./tests/...

coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

run: ## Run the server locally
	@echo "Starting server..."
	DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=booking_service DB_SSLMODE=disable PORT=8080 \
	go run cmd/server/main.go

build: ## Build the application
	@echo "Building application..."
	go build -o bin/booking-service cmd/server/main.go

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/ coverage.out coverage.html

docker-up: ## Start PostgreSQL with docker-compose
	@echo "Starting PostgreSQL..."
	docker run -d --name booking-postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=booking_service \
		-p 5432:5432 \
		postgres:16

docker-down: ## Stop PostgreSQL
	@echo "Stopping PostgreSQL..."
	docker stop booking-postgres || true
	docker rm booking-postgres || true

migrate: ## Run database migrations
	@echo "Running migrations..."
	@for file in $$(ls internal/infrastructure/migrations/*.sql | sort); do \
		echo "Running migration: $$file"; \
		psql -h localhost -U postgres -d booking_service -f $$file || exit 1; \
	done
	@echo "All migrations completed successfully"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem -run=^# ./tests/...

dashboard-deps: ## Install dashboard generation dependencies
	@echo "Installing dashboard dependencies..."
	pip install -r dashboards/requirements.txt

dashboard: ## Generate Grafana dashboard from Python code
	@echo "Generating Grafana dashboard..."
	python3 dashboards/booking_service_dashboard.py > dashboards/booking_service.json
	@echo "Dashboard generated: dashboards/booking_service.json"
	@echo "You can import this file into Grafana or use dashboards/booking_service.json"

openapi-validate: ## Validate OpenAPI specification
	@echo "Validating OpenAPI spec..."
	@command -v openapi-generator-cli >/dev/null 2>&1 || { \
		echo "openapi-generator-cli not found. Install with: npm install -g @openapitools/openapi-generator-cli"; \
		exit 1; \
	}
	openapi-generator-cli validate -i openapi.yaml
