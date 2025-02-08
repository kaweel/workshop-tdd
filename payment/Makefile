# Set default TimeZone
TZ ?= UTC

# Start the application
start-app:
	@echo "🚀 Starting application with TimeZone=$(TZ)..."
	go run main.go
# Run all tests with verbose output
test:
	@echo "🧪 Running all tests..."
	go test -v ./...
# Run only integration tests	
integration-test:
	@echo "🔗 Running integration tests..."
	go test -tags=integration_test ./...
# Run only unit tests
unit-test:
	@echo "🛠 Running unit tests..."
	go test -tags=unit_test ./...
# Run all tests and generate coverage report
test-gen-cov:
	@echo "📊 Running tests with coverage..."	
	go test -tags=unit_test,integration_test ./... -coverprofile=coverage.out
# Show function-wise coverage report
watch-cov-func:
	@echo "📊 Viewing function coverage..."
	go tool cover -func=coverage.out
# Open HTML coverage report
watch-cov:
	@echo "🌐 Opening HTML coverage report..."
	go tool cover -html=coverage.out
# Start the database (Docker Compose)
db-up:
	@echo "🐳 Starting database container..."
	docker compose up
# Stop and remove database container
db-down:
	@echo "🛑 Stopping database container..."
	docker compose down
# Run full cycle: Start DB, Run tests, Stop DB
start:
	@echo "🚀 Starting database, running tests, then stopping database..."
	make start
	make db-up