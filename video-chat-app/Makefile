# Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=videochatapp
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the project
all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Run tests
test:
	$(GOTEST) -count=1 -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Cross compile for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Download dependencies
deps:
	$(GOGET) -v -d ./...

# Run go fmt on all go files
fmt:
	$(GOCMD) fmt ./...

# Run go vet on all go files
vet:
	$(GOCMD) vet ./...

.PHONY: all build test clean run build-linux deps fmt vet