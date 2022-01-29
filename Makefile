all: fmt vet mod lint gen

# Run tests
test: fmt vet
	go test ./...

# Run go fmt against code
fmt:
	go fmt ./...

# Run go fmt against code
mod:
	go mod tidy && go mod verify && \
	buf mod update api

# Run go vet against code
vet:
	go vet ./...

# Run linters
lint:
	buf lint && \
	golangci-lint run

# Run code generators
gen:
	(cd api && buf generate)
