.PHONY: build test clean install uninstall cleaninstall

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
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

uninstall:
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

cleaninstall:
	make uninstall
	make clean
	make build
	make install

# Dependencies
deps:
	$(GOGET) github.com/sashabaranov/go-openai