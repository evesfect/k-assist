.PHONY: build test clean install

# Binary name
BINARY_NAME=kass

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Build directory
BUILD_DIR=build

all: test build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/kass

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

install:
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Dependencies
deps:
	$(GOGET) github.com/sashabaranov/go-openai