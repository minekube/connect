all: fmt vet mod lint

# Run tests
test: fmt vet
	go test ./...

# Run go fmt against code
fmt:
	go fmt ./...

# Run go fmt against code
mod:
	go mod tidy && go mod verify

# Run go vet against code
vet:
	go vet ./...

# Run linters
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

# Serve the docs website locally and auto on changes
dev-docs:
	cd .web && yarn install && yarn dev
