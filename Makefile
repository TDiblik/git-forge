GO=go
BINARY_NAME = git-forge
BINARY_DIR = ./bin
INSTALL_PATH ?= $(HOME)/.local/bin

.PHONY: all build build-prod install update clean use

all: install build

install:
	$(GO) install

build:
	$(GO) build -o $(BINARY_DIR)/$(BINARY_NAME) main.go

build-prod:
	CGO_ENABLED=0 $(GO) build -ldflags="-s -w" -o $(BINARY_DIR)/$(BINARY_NAME) main.go

update: install
	$(GO) mod tidy
	gofmt -w -l .

clean:
	rm -rf $(BINARY_DIR)

use: install build-prod
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Done! Try running 'git forge'."
	@echo "If it does not work, make sure $(INSTALL_PATH) is in your PATH"

clean-use: clean
	rm -rf $(INSTALL_PATH)/$(BINARY_NAME)
