MAIN_BINARY = ./bin/azure-storage-explorer

export GOFLAGS=-mod=vendor

# allows to specify which tests to be run (ex: TEST_PATTERN=FooTest make test)
TEST_PATTERN ?= .

OUTPUT_DIR ?= .

# allow passing -ldflags, etc for release builds
BUILD_ARGS ?=

.PHONY: clean
clean: ## clean up previous builds
	$(call display_msg,Cleaning previous builds)
	@rm -f $(MAIN_BINARY)

.PHONY: build
build: clean ## build all binaries
	@printf ">> building binaries..."
	@go build $(BUILD_ARGS) -o $(MAIN_BINARY) ./cmd/azure-storage-explorer/...
	$(call display_check)

.DEFAULT_GOAL := build
default: build

define display_check
	@printf " \033[32m✔︎\033[0m\n"
endef
