# Project variables
BINARY_NAME=gizmo
BIN_DIR=bin
VERSION?=0.1.0
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

.PHONY: all
all: help

.PHONY: build
## build: Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BIN_DIR}
	$(GOBUILD) $(LDFLAGS) -o ${BIN_DIR}/${BINARY_NAME} .

.PHONY: clean
## clean: Remove binary and build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf ${BIN_DIR}

.PHONY: test
## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: fmt
## fmt: Format source code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

.PHONY: vet
## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

.PHONY: run
## run: Build and run the binary
run: build
	./${BIN_DIR}/${BINARY_NAME}

.PHONY: help
## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '##' $(MAKEFILE_LIST) | grep -v grep | sed -e 's/## //'
