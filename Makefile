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

# Cross-platform build targets
build-linux:
	GOOS=linux GOARCH=amd64 go build -o taskctl-linux-amd64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o taskctl-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o taskctl-darwin-arm64 .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o taskctl-windows-amd64.exe .
