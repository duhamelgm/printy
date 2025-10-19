# Variables
GO := /usr/local/go/bin/go
BINARY := printy
SRC := .

# Cross-compilation settings
GOOS := linux
GOARCH := arm
GOARM := 6

.PHONY: all build build-local clean

# Default target
all: build

# Build for local development (macOS)
build-local:
	@echo "Building $(BINARY) for local development..."
	CGO_ENABLED=1 $(GO) build -o $(BINARY) $(SRC)

# Build for ARM Linux (Raspberry Pi) - requires Docker or cross-compilation tools
build:
	@echo "Building $(BINARY) for $(GOOS)/$(GOARCH) ARMv$(GOARM)..."
	@echo "Note: This requires cross-compilation tools or Docker for ARM Linux"
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) CC=arm-linux-gnueabihf-gcc $(GO) build -o $(BINARY) $(SRC)

# Build using Docker (alternative method)
build-docker:
	@echo "Building $(BINARY) for ARM Linux using Docker..."
	docker run --rm -v "$(PWD)":/usr/src/myapp -w /usr/src/myapp golang:1.23 sh -c "apt-get update && apt-get install -y gcc-arm-linux-gnueabihf && CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 CC=arm-linux-gnueabihf-gcc go build -o $(BINARY) $(SRC)"

# Clean target
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
