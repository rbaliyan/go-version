# Default recipe
default:
    @just --list

# Build all packages
build:
    go build ./...

# Build the CLI binary with version info
build-cli:
    go build -ldflags="$(go run ./cmd/go-version ldflags -static)" -o bin/go-version ./cmd/go-version

# Install the CLI locally
install:
    go install -ldflags="$(go run ./cmd/go-version ldflags -static)" ./cmd/go-version

# Run tests
test:
    go test ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# Run tests with race detector
test-race:
    go test -race ./...

# Run tests with coverage
test-cover:
    go test -cover ./...

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    golangci-lint run ./...

# Tidy dependencies
tidy:
    go mod tidy

# Run vulnerability check
vulncheck:
    go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Check for outdated dependencies
depcheck:
    go list -m -u all | grep '\[' || echo "All dependencies are up to date"

# Build snapshot release (for testing)
snapshot:
    goreleaser release --snapshot --clean

# Check goreleaser config
check-release:
    goreleaser check

# Create and push a new release tag (bumps patch version)
release:
    ./scripts/release.sh

# Clean build artifacts
clean:
    rm -rf bin/ dist/
    go clean -testcache
