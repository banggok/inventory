# Define variables for database settings, Go flags, and test coverage
DB_HOST=localhost
DB_PORT=5432
DB_USER=yourusername
DB_PASSWORD=yourpassword
DB_NAME=inventory_db
MIGRATION_PATH=migrations
GOCMD=go
GINKGO=ginkgo
COVERAGE_OUT=coverage.out

# Default target: Run the service
.PHONY: run
run:
	@echo "Running the Inventory Management service..."
	$(GOCMD) run ./cmd/api

# Migration commands: Run migrations using soda (change soda to your migration tool if different)
.PHONY: migrate
migrate:
	@echo "Running migrations..."
	soda migrate

.PHONY: rollback
rollback:
	@echo "Rolling back migrations..."
	soda rollback

# Run the E2E tests with coverage
.PHONY: test
test:
	@echo "Running E2E tests with coverage..."
	$(GINKGO) -r -cover -coverpkg=./internal/... -coverprofile=$(COVERAGE_OUT)

# Open the coverage report in HTML format
.PHONY: coverage
coverage:
	@echo "Generating coverage report..."
	go tool cover -html=$(COVERAGE_OUT) -o coverage.html
	@echo "Opening coverage report..."
	open coverage.html # On Linux, you might need to use xdg-open

# Clean up: Remove the coverage report file
.PHONY: clean
clean:
	@echo "Cleaning up coverage files..."
	rm -f $(COVERAGE_OUT) coverage.html

# Run all: Run migrations, then run the service and tests
.PHONY: all
all: migrate run test coverage
