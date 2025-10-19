# Variables
GO := /usr/local/go/bin/go
BINARY := printy
SRC := .

# Cross-compilation settings
GOOS := linux
GOARCH := arm
GOARM := 6

.PHONY: all build clean

# Default target
all: build

# Build target
build:
	@echo "Building $(BINARY) for $(GOOS)/$(GOARCH) ARMv$(GOARM)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) build -o $(BINARY) $(SRC)

# Clean target
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
