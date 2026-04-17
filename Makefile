GO=go
BINARY_NAME=git-forge
BINARY_DIR=./bin

.PHONY: all build build-prod install update clean

all: install build

install:
	$(GO) install

build:
	$(GO) build -o $(BINARY_DIR)/$(BINARY_NAME) main.go

build-prod:
	CGO_ENABLED=0 $(GO) build -ldflags="-s -w" -o $(BINARY_DIR)/$(BINARY_NAME) main.go

update:
	$(GO) mod tidy
	gofmt -w -l .

clean:
	rm -rf $(BINARY_DIR)
