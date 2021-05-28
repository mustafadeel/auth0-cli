.DEFAULT_GOAL := help

GOPRIVATE := github.com/auth0

TEST_PATH ?= "./..."
test: ## Run gotestsum and go test on all Go packages.
	@command -v gotestsum >/dev/null && \
		CGO_ENABLED=1 gotestsum --junitfile junit.xml -- $(TEST_PATH) ${tags} ${parallel} -race -count 1 || \
		CGO_ENABLED=1 go test $(TEST_PATH) ${tags} ${parallel} -race -count 1

start: build
	GOREMAN_RPC_ADDR=127.0.0.1 ./.tools/bin/goreman -set-ports=false start

build:
	GOBIN=$$(pwd)/.tools/bin go install ./cmd/logger0

lint: ## Run golangci-lint.
	golangci-lint run -v --timeout=5m

install-tools: ## Install base tooling (protoc, protoc Go plugins, mockgen, etc.)
	./scripts/install_tools.sh

generate: install-tools ## Run the code generation script (protoc, go generate, etc.)
	./scripts/generate.sh

modules: ## Tidy and re-vendor Go modules to ensure dependencies are in a pristine state
	GOPRIVATE=$(GOPRIVATE) go mod tidy
	GOPRIVATE=$(GOPRIVATE) go mod vendor

.PHONY: clean
clean: ## remove files created during build process
	rm -rf .tools

.PHONY: help
help:
	@grep -E '^[%a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
