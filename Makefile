.PHONY: all build-frontend build-go clean run dev

# Default target
all: build-frontend build-go

# Build frontend
build-frontend:
	@echo "Building frontend..."
	cd web/frontend && npm run build

# Build Go binary
build-go:
	@echo "Building Go binary..."
	go build -o taskctl .

# Build everything
build: build-frontend build-go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f taskctl
	rm -rf web/frontend/dist

# Run in development mode (dev server for frontend)
dev:
	@echo "Starting development servers..."
	@echo "Frontend: http://localhost:5173"
	cd web/frontend && npm run dev

# Run the CLI
run:
	go run .

# Install frontend dependencies
install-frontend:
	cd web/frontend && npm install

# Goreleaser targets
release-snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release
